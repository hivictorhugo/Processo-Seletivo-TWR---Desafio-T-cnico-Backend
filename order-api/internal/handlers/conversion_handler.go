package handlers

import (
    "database/sql"
    "encoding/json"
    "io"
    "log"
    "net/http"

    "github.com/username/affiliate-conversions/internal/models"
    "github.com/username/affiliate-conversions/internal/utils"
)

type ConversionHandler struct {
    db *sql.DB
}

func NewConversionHandler(db *sql.DB) *ConversionHandler {
    return &ConversionHandler{db: db}
}

func (h *ConversionHandler) HandleConversions(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    partnerID := r.Header.Get("X-Partner-Id")
    if partnerID == "" {
        http.Error(w, "missing partner id", http.StatusUnauthorized)
        return
    }
    sig := r.Header.Get("X-Signature")
    if sig == "" {
        http.Error(w, "missing signature", http.StatusUnauthorized)
        return
    }

    // fetch partner secret from db
    var secret string
    var partnerName string
    err = h.db.QueryRow(`SELECT secret_key, name FROM partners WHERE partner_id = ?`, partnerID).Scan(&secret, &partnerName)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "unknown partner", http.StatusUnauthorized)
            return
        }
        log.Printf("db partner lookup error: %v", err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    if !utils.ValidateHMAC(body, sig, secret) {
        http.Error(w, "invalid signature", http.StatusUnauthorized)
        return
    }

    var conv models.Conversion
    if err := json.Unmarshal(body, &conv); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }

    // ensure partner id and name match
    conv.PartnerID = partnerID
    if conv.PartnerName == "" {
        conv.PartnerName = partnerName
    }

    // attempt insert
    _, err = h.db.Exec(`INSERT INTO conversions (transaction_id, partner_id, partner_name, sale_amount) VALUES (?, ?, ?, ?)`, conv.TransactionID, conv.PartnerID, conv.PartnerName, conv.SaleAmount)
    if err != nil {
        // check duplicate entry
        if mysqlErr, ok := err.(*mysql.MySQLError); ok {
            // 1062 is duplicate key
            if mysqlErr.Number == 1062 {
                w.WriteHeader(http.StatusOK)
                w.Write([]byte(`{"status":"duplicate"}`))
                return
            }
        }
        log.Printf("db insert error: %v", err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(`{"status":"created"}`))
}