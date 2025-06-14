package models

import (
	"time"
)

type ListingDB struct {
	ID           int64
	Name         string
	Typel        string
	Description  string
	Status       string
	Price        float64
	City         string
	UserID       int64
	Date_created time.Time
}

type Listing struct {
	Name        string `json:"title"`
	Typel       string `json:"type"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Price       int64  `json:"price"`
	City        string `json:"city"`
	UserID      int64
}
