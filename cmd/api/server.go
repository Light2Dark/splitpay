package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/sashabaranov/go-openai"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type application struct {
	logger *slog.Logger
	db     *sql.DB
	openai *openai.Client
	env    string
}

func main() {
	port := flag.Int("port", 8080, "Port number to serve the server")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		slog.Info(".env file not detected")
	}

	var env string = "PROD"
	var logLevel = slog.LevelInfo
	if os.Getenv("ENVIRONMENT") == "DEV" {
		env = "DEV"
		logLevel = slog.LevelDebug
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      logLevel,
		TimeFormat: time.Kitchen,
	}))

	turso_db_url, exists := os.LookupEnv("TUSRO_DB_URL")
	if !exists {
		logger.Warn("unable to obtain turso db url")
	}
	turso_token, exists := os.LookupEnv("TURSO_TOKEN")
	if !exists {
		logger.Warn("unable to obtain turso db token")
	}

	var db_url string
	// TODO: libsql lib does not work without CGO
	// if env == "PROD" {
	// 	db_url = fmt.Sprintf("libsql://%s.turso.io?authToken=%s", turso_db_url, turso_token)
	// } else {
	// 	db_url = "file:./local-sqlite.db"
	// }

	db_url = fmt.Sprintf("libsql://%s.turso.io?authToken=%s", turso_db_url, turso_token)
	db, err := sql.Open("libsql", db_url)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to db %s", err))
	}
	err = db.Ping()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to ping db %s", err))
	}
	defer db.Close()

	openai_token, exists := os.LookupEnv("OPENAI_TOKEN")
	if !exists {
		logger.Error("unable to obtain openai token")
	}
	openai_client := openai.NewClient(openai_token)

	app := application{
		logger: logger,
		db:     db,
		openai: openai_client,
		env:    env,
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Path[len("/static/"):]
		fullPath := filepath.Join(".", "static", filePath)
		http.ServeFile(w, r, fullPath)
	})

	router.HandleFunc("GET /healthcheck", app.healthcheckHandler)
	router.HandleFunc("GET /", app.indexHandler)
	router.HandleFunc("POST /scanReceipt", app.scanReceiptHandler)
	router.HandleFunc("PUT /saveReceipt", app.saveReceiptHandler)
	router.HandleFunc("GET /viewReceipt/{receiptLink}", app.viewReceiptHandler)
	router.HandleFunc("POST /payReceipt", app.payReceiptHandler)
	router.HandleFunc("POST /markPaid", app.markPaidHandler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}

	app.logger.Info(fmt.Sprintf("Starting server on port %s", server.Addr))
	server.ListenAndServe()
}
