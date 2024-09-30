package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
		Help: "Total number of bids placed in all auctions",
	})

	totalValidBids = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auction_total_valid_bids",
		Help: "Total valid number of bids placed in all auctions",
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

const port = 6003

// Estruturas do Bid e Auction
type Bid struct {
	ID     string    `json:"id"`
	Amount float64   `json:"amount"`
	TS     time.Time `json:"timestamp"`
}

type Auction struct {
	ID           string    `json:"id"`
	TimeStart    time.Time `json:"time_start"`
	TimeEnd      time.Time `json:"time_end"`
	MinimumValue float64   `json:"minimum_value"`
	MaximumValue *float64  `json:"maximum_value,omitempty"`
	Bids         []Bid     `json:"bids"`
}

// AuctionHandler gerencia a lógica do leilão e os endpoints
type AuctionHandler struct {
	auctions    map[string]*Auction
	subscribers map[string][]*websocket.Conn // Subscrições dos WebSockets por leilão
	mu          sync.Mutex
	redisClient *redis.Client // Cliente Redis
}

// Novo AuctionHandler para inicializar os leilões
func NewAuctionHandler(redisClient *redis.Client) *AuctionHandler {
	return &AuctionHandler{
		auctions:    make(map[string]*Auction),
		subscribers: make(map[string][]*websocket.Conn),
		redisClient: redisClient,
	}
}

// Adiciona um leilão ao mapa de leilões
func (ah *AuctionHandler) AddAuction(a *Auction) {
	ah.auctions[a.ID] = a
}

// Valida se o leilão está ativo no momento
func (a *Auction) isAuctionActive() bool {
	now := time.Now()
	return now.After(a.TimeStart) && now.Before(a.TimeEnd)
}

// Valida um lance baseado nas regras do leilão
func (a *Auction) isValidBid(bid Bid) error {
	if bid.Amount < a.MinimumValue {
		return fmt.Errorf("bid amount too low")
	}
	if a.MaximumValue != nil && bid.Amount > *a.MaximumValue {
		return fmt.Errorf("bid exceeds maximum allowed value")
	}
	return nil
}

// Adiciona o lance ao leilão e notifica os clientes WebSocket
func (ah *AuctionHandler) addBidAndNotify(auction *Auction, bid Bid) {
	bid.TS = time.Now()
	auction.Bids = append(auction.Bids, bid)

	bidAmountHistogram.Observe(bid.Amount)
	totalValidBids.Inc()

	// Notificar os clientes WebSocket conectados
	ah.notifySubscribers(auction.ID, bid)
}

// Notifica os subscritores WebSocket de novos lances
func (ah *AuctionHandler) notifySubscribers(auctionID string, bid Bid) {
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

func (ah *AuctionHandler) acquireLock(ctx context.Context, auctionID string) (bool, error) {
	lockKey := fmt.Sprintf("auction_lock:%s", auctionID)
	success, err := ah.redisClient.SetNX(ctx, lockKey, 1, 10*time.Second).Result() // Expira em 10 segundos
	return success, err
}

func (ah *AuctionHandler) releaseLock(ctx context.Context, auctionID string) error {
	lockKey := fmt.Sprintf("auction_lock:%s", auctionID)
	_, err := ah.redisClient.Del(ctx, lockKey).Result()
	return err
}

// HandleBid recebe um POST request para lances
func (ah *AuctionHandler) HandleBid(ctx *fasthttp.RequestCtx) {
	totalBids.Inc()

	auctionID := ctx.UserValue("auction_id").(string)

	auction, ok := ah.auctions[auctionID]
	if !ok {
		ctx.Error("Auction not found", fasthttp.StatusNotFound)
		return
	}

	var bid Bid
	if err := json.Unmarshal(ctx.PostBody(), &bid); err != nil {
		ctx.Error("Invalid bid format", fasthttp.StatusBadRequest)
		return
	}

	// Adquire lock no Redis
	ctxRedis := context.Background()
	if acquired, err := ah.acquireLock(ctxRedis, auctionID); !acquired || err != nil {
		ctx.Error("Could not acquire lock", fasthttp.StatusInternalServerError)
		return
	}
	defer ah.releaseLock(ctxRedis, auctionID) // Libera o lock no final

	// Validações
	if !auction.isAuctionActive() {
		ctx.Error("Auction is not active", fasthttp.StatusForbidden)
		return
	}

	if err := auction.isValidBid(bid); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusForbidden)
		return
	}

	// Adiciona o lance e notifica os clientes conectados via WebSocket
	ah.addBidAndNotify(auction, bid)
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
	auctionID := ctx.UserValue("auction_id").(string)

	auction, ok := ah.auctions[auctionID]
	if !ok {
		ctx.Error("Auction not found", fasthttp.StatusNotFound)
		return
	}

	// Estabelece a conexão WebSocket
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
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

	auctionHandler := NewAuctionHandler(redisClient)

	// Exemplo de leilão ativo por 5 minutos
	auctionID := "auction1"
	auctionHandler.AddAuction(&Auction{
		ID:           auctionID,
		TimeStart:    time.Now(),
		TimeEnd:      time.Now().Add(10 * time.Minute),
		MinimumValue: 10.0,
		MaximumValue: nil,
		Bids:         []Bid{},
	})

	// Configura o roteador
	r := router.New()
	r.POST("/auctions/{auction_id}/bids", auctionHandler.HandleBid)
	r.GET("/auctions/{auction_id}/bids", auctionHandler.handleWebSocket)

	p := fastp.NewPrometheus("fasthttp")
	fastpHandler := p.WrapHandler(r)

	// Inicializa o servidor
	fmt.Println("Auction active server is running on port 6003...")
	if err := fasthttp.ListenAndServe(":6003", fastpHandler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
