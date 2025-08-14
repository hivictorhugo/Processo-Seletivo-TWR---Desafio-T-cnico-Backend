package main

import (
    "log"
    "net/http"
    "os"

    _ "github.com/go-sql-driver/mysql"
    "github.com/username/affiliate-conversions/internal/db"
    "github.com/username/affiliate-conversions/internal/handlers"
    "github.com/username/affiliate-conversions/internal/middleware"
)

func main() {
    dsn := os.Getenv("DB_DSN")
    if dsn == "" {
        log.Fatal("DB_DSN is required")
    }

    database, err := db.New(dsn)
    if err != nil {
        log.Fatalf("db connect: %v", err)
    }
    defer database.Close()

    h := handlers.NewConversionHandler(database)

    mux := http.NewServeMux()
    mux.Handle("/conversions", middleware.Logging(middleware.LimitBody(http.HandlerFunc(h.HandleConversions), 1<<20)))

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("listening on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, mux))
}