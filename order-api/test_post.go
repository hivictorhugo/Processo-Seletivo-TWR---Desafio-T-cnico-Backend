package main

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "net/http"
)

func main() {
    url := "http://localhost:8080/conversions"
    body := []byte(`{"transaction_id":"tx_test_1","partner_name":"Partner A","sale_amount":99.90}`)
    secret := "secret_for_partner_a"

    // gerar HMAC
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(body)
    sig := hex.EncodeToString(mac.Sum(nil))

    req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
    req.Header.Set("X-Partner-Id", "partner-a")
    req.Header.Set("X-Signature", sig)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    respBody, _ := io.ReadAll(resp.Body)
    fmt.Println("Status:", resp.Status)
    fmt.Println("Response:", string(respBody))
}
