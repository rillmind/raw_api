package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rillmind/raw_api/db"
	"github.com/rillmind/raw_api/src/product"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Erro ao carregar variaveis de ambiente: %v", err)
	}

	dsn := os.Getenv("DATA_SOURCE_NAME")
	if dsn == "" {
		dsn = "postgres://user:123456@localhost:5432/raw_api_db?sslmode=disable"
	}

	db, err := db.New(dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar ao bd: %v", err)
	}

	defer db.Close()

	store := product.NewRepository(db)
	storeHandler := product.NewStoreHandler(*store)

	mux := http.NewServeMux()
	storeHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":2306",
		Handler:      mux,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}

	go func() {
		log.Println("Servidor rodando em http://127.0.0.1:2306")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro no servidor: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Println("Desligando servidor...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor: %v", err)
	}
}
