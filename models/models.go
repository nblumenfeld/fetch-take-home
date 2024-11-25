package models

import (
	"github.com/google/uuid"
)

// receipt object and the in memory cache
type Receipt struct {
	ID           uuid.UUID `json:"id"`
	Retailer     string    `json:"retailer"`
	PurchaseDate string    `json:"purchaseDate"`
	PurchaseTime string    `json:"purchaseTime"`
	Total        string    `json:"total"`
	Items        []Item    `json:"items"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Points struct {
	ID     uuid.UUID
	Points int
}
