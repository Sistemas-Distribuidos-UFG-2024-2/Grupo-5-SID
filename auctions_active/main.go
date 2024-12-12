package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

// Estruturas do Bid e AuctionActive
type Bid struct {
	ID           string    `json:"id"`
	AccountEmail string    `json:"account_email"`
	Amount       float64   `json:"amount"`
	TS           time.Time `json:"timestamp"`
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

	subscribers    map[int64][]*websocket.Conn // Subscrições dos WebSockets por leilão
	activeChannels map[int64]chan struct{}     // Canais ativos para escuta do Redis
	mu             sync.Mutex
	redisClient    *redis.Client // Cliente Redis
}

func NewAuctionHandler(redisClient *redis.Client, auctionClient *AuctionClient) *AuctionHandler {
	return &AuctionHandler{
		auctionClient:  auctionClient,
		subscribers:    make(map[int64][]*websocket.Conn),
		activeChannels: make(map[int64]chan struct{}),
		redisClient:    redisClient,
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

	channel := fmt.Sprintf("auction:%d:channel", auction.ID)
	if err := ah.redisClient.Publish(context.Background(), channel, bidData).Err(); err != nil {
		log.Print("Could not publish bid")
		return
	}

	// Notificar os clientes WebSocket conectados
	// ah.notifySubscribers(auction.ID, bid)
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
			ctx.Error("Auction not found: "+err.Error(), fasthttp.StatusNotFound)
			return
		}

		if !auction.isActive() {
			ctx.Error("AuctionActive not active", fasthttp.StatusNoContent)
			return
		}

		auctionActive, err = auction.convertToActive()
		if err != nil {
			ctx.Error("Could not convert to active auction:"+err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		err = ah.CreateOrUpdate(auctionActive)
		if err != nil {
			ctx.Error("Could not add auction to Redis:"+err.Error(), fasthttp.StatusInternalServerError)
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

	go func() {
		if err := ah.auctionClient.Auction(auctionID, bid); err != nil {
			fmt.Printf("error to update the auction service: %v", err)
		}
	}()

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

func (ah *AuctionHandler) subscribeToRedis(auctionID int64) {
	if _, exists := ah.activeChannels[auctionID]; exists {
		return
	}

	stopChan := make(chan struct{})
	ah.activeChannels[auctionID] = stopChan

	go func() {
		ctx := context.Background()
		channel := fmt.Sprintf("auction:%d:channel", auctionID)
		pubsub := ah.redisClient.Subscribe(ctx, channel)
		if _, err := pubsub.Receive(ctx); err != nil {
			log.Printf("Error subscribing to Redis channel: %v", err)
			return
		}

		log.Printf("Subscribed to Redis channel: %s", channel)

		defer pubsub.Close()

		for msg := range pubsub.Channel() {
			log.Printf("Message received on channel %s: %s", channel, msg.Payload)

			var bid Bid
			if err := json.Unmarshal([]byte(msg.Payload), &bid); err != nil {
				log.Printf("failed to unmarshal bid message: %v", err)
				continue
			}

			conns := ah.subscribers[auctionID]

			for _, conn := range conns {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
					log.Printf("Error sending bid to subscriber: %v", err)
					conn.Close() // Fechar a conexão com erro
				}
			}
		}
	}()
}

func (ah *AuctionHandler) stopRedisSubscription(auctionID int64) {
	ah.mu.Lock()
	if stopChan, exists := ah.activeChannels[auctionID]; exists {
		close(stopChan)
		delete(ah.activeChannels, auctionID)
	}
	ah.mu.Unlock()
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
		if len(ah.subscribers[auctionID]) == 1 {
			ah.subscribeToRedis(auctionID) // Inicia a escuta Redis se for o primeiro cliente
		}
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

		ah.mu.Lock()
		for i, subscriber := range ah.subscribers[auctionID] {
			if subscriber == conn {
				ah.subscribers[auctionID] = append(ah.subscribers[auctionID][:i], ah.subscribers[auctionID][i+1:]...)
				break
			}
		}
		if len(ah.subscribers[auctionID]) == 0 {
			ah.stopRedisSubscription(auctionID) // Para a escuta Redis se não houver mais clientes
		}
		ah.mu.Unlock()
	})
	if err != nil {
		log.Printf("Failed to establish WebSocket connection: %v", err)
	}
}

// Vou rodar no K8S ?
const K8S = true

func main() {
	var (
		redisBaseURL         = "localhost:6379"
		auctionClientBaseURL = "http://localhost:8080"
	)

	if K8S {
		redisBaseURL = "redis-service:6379"
		// auctionClientBaseURL = "http://local-service:8080"
		auctionClientBaseURL = "http://host.docker.internal:8080"
	}

	port := "6003"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	if len(os.Args) > 2 {
		arg := os.Args[2]
		if arg == "compose" {
			println("COMPOSE")
			auctionClientBaseURL = "http://springboot-app:8080"
			redisBaseURL = "redis:6379"
		}
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisBaseURL, // Endereço do servidor Redis
	})

	client := &AuctionClient{auctionClientBaseURL, 1}
	auctionHandler := NewAuctionHandler(redisClient, client)

	// Configura o roteador
	r := router.New()
	r.POST("/auctions/{auction_id}/bids", auctionHandler.HandleBid)
	r.GET("/auctions/{auction_id}/bids/ws", auctionHandler.handleWebSocket)

	p := fastp.NewPrometheus("fasthttp")
	fastpHandler := p.WrapHandler(r)

	handlerWithCors := func(ctx *fasthttp.RequestCtx) {
		corsMiddleware(ctx)
		fastpHandler(ctx)
	}

	// Inicializa o servidor
	fmt.Println("Auction active server is running on port " + port)
	if err := fasthttp.ListenAndServe(":"+port, handlerWithCors); err != nil {
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

func corsMiddleware(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")                                // Permite todas as origens
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Métodos permitidos
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")                               // Permite todos os cabeçalhos
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")                        // Permite cookies/credenciais

	// Se for uma requisição OPTIONS, apenas responde com os cabeçalhos e finaliza a requisição
	if string(ctx.Method()) == "OPTIONS" {
		ctx.SetStatusCode(fasthttp.StatusNoContent) // Responde sem conteúdo (204)
		return
	}

	// Continuar para o próximo handler
}
