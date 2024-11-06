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
	ID int64 `json:"id"`
}

func (a Auction) isActive() bool {
	return true
}

func (a Auction) convertToActive() (*AuctionActive, error) {
	// if not active returns error

	return &AuctionActive{
		ID:           a.ID,
		TimeStart:    time.Now(),
		TimeEnd:      time.Now().Add(10 * time.Minute),
		MinimumValue: 10.0,
		MaximumValue: nil,
		Bids:         []Bid{},
	}, nil
}

func NewAuctionClient(baseURL string) *AuctionClient {
	return &AuctionClient{BaseURL: baseURL}
}

func (c *AuctionClient) GetByID(id int64) (Auction, error) {
	if c.mockID != 0 {
		return Auction{ID: c.mockID}, nil
	}

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

	return auction, nil
}
