package models

import "time"

type Partner struct {
    PartnerID string
    Name      string
    SecretKey string
    CreatedAt time.Time
}

type Conversion struct {
    ID            int64   `json:"-"`
    TransactionID string  `json:"transaction_id"`
    PartnerID     string  `json:"partner_id"`
    PartnerName   string  `json:"partner_name"`
    SaleAmount    float64 `json:"sale_amount"`
    CreatedAt     time.Time `json:"created_at"`
}