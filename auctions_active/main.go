package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/fasthttp/router"
	"github.com/fasthttp/websocket"
	fastp "github.com/flf2ko/fasthttp-prometheus"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"
)

var (
	totalBids = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auction_total_bids",
		Help: "Total number of bids placed in all auctionsActive",
	})

	totalValidBids = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auction_total_valid_bids",
		Help: "Total valid number of bids placed in all auctionsActive",
	})

	bidAmountHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "auction_bid_amounts",
		Help:    "Histogram of bid amounts",
		Buckets: prometheus.LinearBuckets(10, 10, 10), // Intervalos de 10 para o histograma
	})
)

func init() {
	// Registrando o contador com o Prometheus
	prometheus.MustRegister(totalBids)
	prometheus.MustRegister(totalValidBids)
	prometheus.MustRegister(bidAmountHistogram)
}

const port = "6003"

// Estruturas do Bid e AuctionActive
type Bid struct {
	ID     string    `json:"id"`
	Amount float64   `json:"amount"`
	TS     time.Time `json:"timestamp"`
}

type AuctionActive struct {
	ID           int64     `json:"id"`
	TimeStart    time.Time `json:"time_start"`
	TimeEnd      time.Time `json:"time_end"`
	MinimumValue float64   `json:"minimum_value"`
	MaximumValue *float64  `json:"maximum_value,omitempty"`
	Bids         []Bid     `json:"bids"`
}

// AuctionHandler gerencia a lógica do leilão e os endpoints
type AuctionHandler struct {
	auctionClient *AuctionClient

	auctionsActive map[int64]*AuctionActive
	subscribers    map[int64][]*websocket.Conn // Subscrições dos WebSockets por leilão
	mu             sync.Mutex
	redisClient    *redis.Client // Cliente Redis
}

func NewAuctionHandler(redisClient *redis.Client, auctionClient *AuctionClient) *AuctionHandler {
	return &AuctionHandler{
		auctionClient:  auctionClient,
		auctionsActive: make(map[int64]*AuctionActive),
		subscribers:    make(map[int64][]*websocket.Conn),
		redisClient:    redisClient,
	}
}

func (ah *AuctionHandler) AddAuction(a *AuctionActive) {
	ah.auctionsActive[a.ID] = a
}

func (a *AuctionActive) isActive() bool {
	now := time.Now()
	return now.After(a.TimeStart) && now.Before(a.TimeEnd)
}

func (a *AuctionActive) isValidBid(bid Bid) error {
	if bid.Amount < a.MinimumValue {
		return fmt.Errorf("bid amount too low")
	}
	if a.MaximumValue != nil && bid.Amount > *a.MaximumValue {
		return fmt.Errorf("bid exceeds maximum allowed value")
	}
	return nil
}

// Adiciona o lance ao leilão e notifica os clientes WebSocket
func (ah *AuctionHandler) addBidAndNotify(auction *AuctionActive, bid Bid) {
	bid.TS = time.Now()
	auction.Bids = append(auction.Bids, bid)

	bidAmountHistogram.Observe(bid.Amount)
	totalValidBids.Inc()

	// Notificar os clientes WebSocket conectados
	ah.notifySubscribers(auction.ID, bid)
}

// Notifica os subscritores WebSocket de novos lances
func (ah *AuctionHandler) notifySubscribers(auctionID int64, bid Bid) {
	ah.mu.Lock()
	defer ah.mu.Unlock()

	conns := ah.subscribers[auctionID]
	for _, conn := range conns {
		bidData, _ := json.Marshal(bid)
		if err := conn.WriteMessage(websocket.TextMessage, bidData); err != nil {
			log.Printf("Error sending bid to subscriber: %v", err)
			conn.Close() // Fechar conexão em caso de erro
		}
	}
}

func (ah *AuctionHandler) acquireLock(ctx context.Context, auctionID int64) (bool, error) {
	lockKey := fmt.Sprintf("auction_lock:%d", auctionID)
	success, err := ah.redisClient.SetNX(ctx, lockKey, 1, 10*time.Second).Result() // Expira em 10 segundos
	return success, err
}

func (ah *AuctionHandler) releaseLock(ctx context.Context, auctionID int64) error {
	lockKey := fmt.Sprintf("auction_lock:%d", auctionID)
	_, err := ah.redisClient.Del(ctx, lockKey).Result()
	return err
}

// HandleBid recebe um POST request para lances
func (ah *AuctionHandler) HandleBid(ctx *fasthttp.RequestCtx) {
	totalBids.Inc()

	auctionID, err := strconv.ParseInt(ctx.UserValue("auction_id").(string), 10, 64)
	if err != nil {
		ctx.Error("Invalid bid format", fasthttp.StatusBadRequest)
		return
	}

	auction, err := ah.auctionClient.GetByID(auctionID)
	if err != nil {
		ctx.Error("Auction not found", fasthttp.StatusNotFound)
		return
	}

	auctionActive, ok := ah.auctionsActive[auctionID]
	if !ok {
		if !auction.isActive() {
			ctx.Error("AuctionActive not active", fasthttp.StatusNoContent)
			return
		}

		auctionActive, err = auction.convertToActive()
		if err != nil {
			ctx.Error("Could not convert to active auction", fasthttp.StatusInternalServerError)
			return
		}
		ah.AddAuction(auctionActive)
	}

	var bid Bid
	if err := json.Unmarshal(ctx.PostBody(), &bid); err != nil {
		ctx.Error("Invalid bid format", fasthttp.StatusBadRequest)
		return
	}

	ctxRedis := context.Background()
	if acquired, err := ah.acquireLock(ctxRedis, auctionID); !acquired || err != nil {
		ctx.Error("Could not acquire lock", fasthttp.StatusInternalServerError)
		return
	}
	defer ah.releaseLock(ctxRedis, auctionID) // Libera o lock no final

	if !auctionActive.isActive() {
		ctx.Error("AuctionActive is not active", fasthttp.StatusForbidden)
		return
	}

	if err := auctionActive.isValidBid(bid); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusForbidden)
		return
	}

	// Adiciona o lance e notifica os clientes conectados via WebSocket
	ah.addBidAndNotify(auctionActive, bid)
	ctx.SetStatusCode(fasthttp.StatusAccepted)
	fmt.Fprintf(ctx, "Bid received successfully")
}

// WebSocket para listar lances em tempo real
var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true // Permite qualquer origem (CORS)
	},
}

func (ah *AuctionHandler) handleWebSocket(ctx *fasthttp.RequestCtx) {
	auctionID, err := strconv.ParseInt(ctx.UserValue("auction_id").(string), 10, 64)

	auction, ok := ah.auctionsActive[auctionID]
	if !ok {
		ctx.Error("AuctionActive not found", fasthttp.StatusNotFound)
		return
	}

	// Estabelece a conexão WebSocket
	err = upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		defer conn.Close()

		// Adiciona o WebSocket aos subscritores
		ah.mu.Lock()
		ah.subscribers[auctionID] = append(ah.subscribers[auctionID], conn)
		ah.mu.Unlock()

		// Envia os lances já existentes
		for _, bid := range auction.Bids {
			bidData, _ := json.Marshal(bid)
			if err := conn.WriteMessage(websocket.TextMessage, bidData); err != nil {
				log.Printf("Error sending initial bids: %v", err)
				return
			}
		}

		// Aguarda até que a conexão WebSocket seja fechada
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("WebSocket connection closed: %v", err)
				break
			}
		}
	})
	if err != nil {
		log.Printf("Failed to establish WebSocket connection: %v", err)
	}
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Endereço do servidor Redis
	})

	client := &AuctionClient{"http://localhost:8080/", 1}
	auctionHandler := NewAuctionHandler(redisClient, client)

	// Configura o roteador
	r := router.New()
	r.POST("/auctions/{auction_id}/bids", auctionHandler.HandleBid)
	r.GET("/auctions/{auction_id}/bids", auctionHandler.handleWebSocket)

	p := fastp.NewPrometheus("fasthttp")
	fastpHandler := p.WrapHandler(r)

	// Inicializa o servidor
	fmt.Println("Auction active server is running on port " + port)
	if err := fasthttp.ListenAndServe(":"+port, fastpHandler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
