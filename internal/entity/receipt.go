package entity

import "time"

type ReceiptDeserialize struct {
	Retailer     string             `json:"retailer,omitempty"`
	PurchaseDate string             `json:"purchaseDate,omitempty"`
	PurchaseTime string             `json:"purchaseTime,omitempty"`
	Total        string             `json:"total,omitempty"`
	Items        []*ItemDeserialize `json:"items,omitempty"`
}

type ItemDeserialize struct {
	ShortDescription string `json:"shortDescription,omitempty"`
	Price            string `json:"price,omitempty"`
}

type Receipt struct {
	Retailer     string    `json:"retailer,omitempty"`
	PurchaseDate time.Time `json:"purchaseDate,omitempty"`
	PurchaseTime time.Time `json:"purchaseTime,omitempty"`
	Total        float64   `json:"total,omitempty"`
	Items        []*Item   `json:"items,omitempty"`
}

type Item struct {
	ShortDescription string  `json:"shortDescription,omitempty"`
	Price            float64 `json:"price,omitempty"`
}

type PointsResponse struct {
	Points float64 `json:"points"`
}

type ReceiptResponse struct {
	Id string `json:"id"`
}
