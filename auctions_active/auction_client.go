package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuctionClient struct {
	BaseURL string

	mockID int64
}

type Auction struct {
	ID              int64   `json:"id"`
	Produto         string  `json:"produto"`
	LanceInicial    float64 `json:"lanceInicial"`
	DataFinalizacao string  `json:"dataFinalizacao"`
	Criador         string  `json:"criador"`
	Vencedor        string  `json:"vencedor"`
	ValorMaximo     float64 `json:"valorMaximo"`
}

func (a Auction) isActive() bool {
	return true
}

func (a Auction) convertToActive() (*AuctionActive, error) {
	timeEnd, err := parseISO8601(a.DataFinalizacao)
	if err != nil {
		return nil, err
	}
	var valorMaximo float64
	if a.ValorMaximo != 0 {
		valorMaximo = a.ValorMaximo
	}

	return &AuctionActive{
		ID:           a.ID,
		TimeStart:    time.Now(),
		TimeEnd:      timeEnd,
		MinimumValue: a.LanceInicial,
		MaximumValue: &valorMaximo,
		Bids:         []Bid{},
	}, nil
}

func (c *AuctionClient) GetByID(id int64) (Auction, error) {
	//if c.mockID != 0 {
	//	return Auction{
	//		ID:              id,
	//		Produto:         "Produto Exemplo",
	//		LanceInicial:    50.00,
	//		DataFinalizacao: "2024-12-31T23:59:59",
	//		Criador:         "roberta",
	//		Vencedor:        "",
	//		ValorMaximo:     0,
	//	}, nil
	//}

	url := fmt.Sprintf("%s/auctions/%d", c.BaseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return Auction{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Auction{}, fmt.Errorf("falha ao buscar leilão: status %d", resp.StatusCode)
	}

	var auction Auction
	if err := json.NewDecoder(resp.Body).Decode(&auction); err != nil {
		return Auction{}, err
	}

	fmt.Printf("Auction found: %d, minimum value:%f\n", auction.ID, auction.LanceInicial)

	return auction, nil
}

func parseISO8601(input string) (time.Time, error) {
	const layout = "2006-01-02T15:04:05"
	parsedTime, err := time.Parse(layout, input)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
