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

	server := server.NewServer(cfg.Host, cfg.Port, store)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
