package models

import "time"

// Order is the model representation of an order in the data model.
type Order struct {
	ID           int         `json:"id"`
	UserID       int         `json:"-"`
	RestaurantID int         `json:"-"`
	Restaurant   *Restaurant `json:"restaurant" sql:"-"`
	Total        int         `json:"total"`
	CurrencyCode string      `json:"currency_code"`
	PlacedAt     time.Time   `json:"placed_at"`
}

// Orders is a slice of Order pointers.
type Orders []*Order
