package tests

import (
    "bytes"
    "encoding/hex"
    "net/http"
    "net/http/httptest"
    "testing"
    "crypto/hmac"
    "crypto/sha256"

    _ "github.com/go-sql-driver/mysql"
    "github.com/username/affiliate-conversions/internal/handlers"
    "github.com/username/affiliate-conversions/internal/db"
)

func TestHandleConversions_CreateAndDuplicate(t *testing.T) {

    dsn := "appuser:apppass@tcp(localhost:3306)/affiliate_db?parseTime=true"
    database, err := db.New(dsn)
    if err != nil {
        t.Fatalf("db connect: %v", err)
    }
    defer database.Close()

    h := handlers.NewConversionHandler(database)

    body := []byte(`{"transaction_id":"tx_test_1","partner_name":"Partner A","sale_amount":99.90}`)
    mac := hmac.New(sha256.New, []byte("secret_for_partner_a"))
    mac.Write(body)
    sig := hex.EncodeToString(mac.Sum(nil))

    req := httptest.NewRequest(http.MethodPost, "/conversions", bytes.NewReader(body))
    req.Header.Set("X-Partner-Id", "partner-a")
    req.Header.Set("X-Signature", sig)
    w := httptest.NewRecorder()

    h.HandleConversions(w, req)
    if w.Result().StatusCode != http.StatusCreated && w.Result().StatusCode != http.StatusOK {
        t.Fatalf("expected 201 or 200, got %d", w.Result().StatusCode)
    }

    // second call should be duplicate
    w2 := httptest.NewRecorder()
    h.HandleConversions(w2, req)
    if w2.Result().StatusCode != http.StatusOK {
        t.Fatalf("expected 200 on duplicate, got %d", w2.Result().StatusCode)
    }
}
