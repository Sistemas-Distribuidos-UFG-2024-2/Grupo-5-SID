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
	LastValue    float64   `json:"last_value"`
	Bids         []Bid     `json:"bids"`
}

// AuctionHandler gerencia a lógica do leilão e os endpoints
type AuctionHandler struct {
	auctionClient *AuctionClient

	subscribers map[int64][]*websocket.Conn // Subscrições dos WebSockets por leilão
	mu          sync.Mutex
	redisClient *redis.Client // Cliente Redis
}

func NewAuctionHandler(redisClient *redis.Client, auctionClient *AuctionClient) *AuctionHandler {
	return &AuctionHandler{
		auctionClient: auctionClient,
		subscribers:   make(map[int64][]*websocket.Conn),
		redisClient:   redisClient,
	}
}

func (ah *AuctionHandler) CreateOrUpdate(a *AuctionActive) error {
	data, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("failed to marshal auction: %v", err)
	}

	key := fmt.Sprintf("auction:%d", a.ID)
	err = ah.redisClient.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to save auction to redis: %v", err)
	}

	return nil
}

func (ah *AuctionHandler) GetAuctionActiveByID(id int64) (*AuctionActive, error) {
	key := fmt.Sprintf("auction:%d", id)
	data, err := ah.redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get auction from redis: %v", err)
	}

	var auction AuctionActive
	if err := json.Unmarshal([]byte(data), &auction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal auction data: %v", err)
	}

	return &auction, nil
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
	if bid.Amount <= a.LastValue {
		return fmt.Errorf("bid amount too low")
	}
	return nil
}

func (ah *AuctionHandler) addBidAndNotify(auction *AuctionActive, bid Bid) {
	bid.TS = time.Now()

	bidData, err := json.Marshal(bid)
	if err != nil {
		log.Printf("failed to marshal bid: %v", err)
		return
	}
	bidKey := fmt.Sprintf("auction:%d:bids", auction.ID)
	err = ah.redisClient.RPush(context.Background(), bidKey, bidData).Err()
	if err != nil {
		log.Printf("failed to save bid to redis: %v", err)
		return
	}

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

	auctionActive, err := ah.GetAuctionActiveByID(auctionID)
	if err != nil {
		println("error to get auction active by id:", err.Error())

		auction, err := ah.auctionClient.GetByID(auctionID)
		if err != nil {
			ctx.Error("Auction not found", fasthttp.StatusNotFound)
			return
		}

		if !auction.isActive() {
			ctx.Error("AuctionActive not active", fasthttp.StatusNoContent)
			return
		}

		auctionActive, err = auction.convertToActive()
		if err != nil {
			ctx.Error("Could not convert to active auction", fasthttp.StatusInternalServerError)
			return
		}

		err = ah.CreateOrUpdate(auctionActive)
		if err != nil {
			ctx.Error("Could not add auction to Redis", fasthttp.StatusInternalServerError)
			return
		}
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
		ctx.Error(err.Error(), fasthttp.StatusAccepted)
		return
	}

	auctionActive.LastValue = bid.Amount
	fmt.Println(fmt.Sprintf("Updating the last value to %f on auction %d", auctionActive.LastValue, auctionActive.ID))

	// Adiciona o lance e notifica os clientes conectados via WebSocket
	ah.addBidAndNotify(auctionActive, bid)
	if err := ah.CreateOrUpdate(auctionActive); err != nil {
		fmt.Println("Error updating auction:", err)
	}
	ctx.SetStatusCode(fasthttp.StatusAccepted)
	fmt.Fprintf(ctx, fmt.Sprintf("Bid accepted: %v", bid))
}

// WebSocket para listar lances em tempo real
var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true // Permite qualquer origem (CORS)
	},
}

func (ah *AuctionHandler) handleWebSocket(ctx *fasthttp.RequestCtx) {
	auctionID, err := strconv.ParseInt(ctx.UserValue("auction_id").(string), 10, 64)
	if err != nil {
		ctx.Error("Invalid auction ID", fasthttp.StatusBadRequest)
		return
	}

	_, err = ah.GetAuctionActiveByID(auctionID)
	if err != nil {
		ctx.Error("Auction not found", fasthttp.StatusNotFound)
		return
	}

	// Estabelece a conexão WebSocket
	err = upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		defer conn.Close()

		// Adiciona o WebSocket aos subscritores
		ah.mu.Lock()
		ah.subscribers[auctionID] = append(ah.subscribers[auctionID], conn)
		ah.mu.Unlock()

		// Recupera e envia os lances do Redis
		bidKey := fmt.Sprintf("auction:%d:bids", auctionID)
		bidData, err := ah.redisClient.LRange(context.Background(), bidKey, 0, -1).Result()
		if err != nil {
			log.Printf("failed to retrieve bids from redis: %v", err)
			return
		}

		for _, bidJSON := range bidData {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(bidJSON)); err != nil {
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

// K8s: redis-service
// Local: localhost
var redisBaseURL = "redis-service"

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis-service:6379", // Endereço do servidor Redis
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

func (ah *AuctionHandler) GetBidsByAuctionID(auctionID int64) ([]Bid, error) {
	// Definimos a chave que armazena os lances do leilão específico
	bidKey := fmt.Sprintf("auction:%d:bids", auctionID)

	// Busca todos os lances da lista
	bidData, err := ah.redisClient.LRange(context.Background(), bidKey, 0, -1).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("no bids found for auction %d", auctionID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve bids from redis: %v", err)
	}

	// Convertendo os dados JSON recuperados em structs Bid
	var bids []Bid
	for _, bidJSON := range bidData {
		var bid Bid
		if err := json.Unmarshal([]byte(bidJSON), &bid); err != nil {
			return nil, fmt.Errorf("failed to unmarshal bid data: %v", err)
		}
		bids = append(bids, bid)
	}

	return bids, nil
}
