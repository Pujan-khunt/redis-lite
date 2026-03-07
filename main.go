package main

import (
	"log"

	"github.com/Pujan-khunt/redis-lite/config"
	"github.com/Pujan-khunt/redis-lite/server"
	"github.com/Pujan-khunt/redis-lite/storage"
)

func main() {
	cfg := config.NewConfig()

	store := storage.NewInMemoryStore()

	server, err := server.NewServer(cfg.Host, cfg.Port, store, cfg.Period)
	if err != nil {
		log.Fatalln("Error creating server", err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("Error starting server", err)
	}
}
