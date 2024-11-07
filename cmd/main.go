package main

import (
	"log"

	"github.com/kkstas/redis-go/internal/config"
	"github.com/kkstas/redis-go/internal/server"
	"github.com/kkstas/redis-go/internal/store"
)

func main() {
	srvConfig := config.New()
	srvStore := store.New()
	srv := server.New(srvConfig, srvStore)

	network := "tcp"
	address := "0.0.0.0:6379"

	log.Println("Starting server at", network, address)
	if err := server.Run(srv, network, address); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
