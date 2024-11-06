package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AccountClient struct {
	BaseURL string

	mockID int64
}

type Account struct {
	ID int64 `json:"id"`
}

func NewAccountClient(baseURL string) *AccountClient {
	return &AccountClient{BaseURL: baseURL}
}

func (c *AccountClient) GetByID(id int64) (Account, error) {
	if c.mockID != 0 {
		return Account{ID: c.mockID}, nil
	}

	url := fmt.Sprintf("%s/accounts/%d", c.BaseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return Account{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Account{}, fmt.Errorf("falha ao buscar conta: status %d", resp.StatusCode)
	}

	var auction Account
	if err := json.NewDecoder(resp.Body).Decode(&auction); err != nil {
		return Account{}, err
	}

	return auction, nil
}
