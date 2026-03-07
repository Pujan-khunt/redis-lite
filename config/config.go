package config

import (
	"flag"
	"time"
)

var (
	port   = flag.Int("port", 6379, "port for redis-lite to listen on")
	host   = flag.String("ip", "0.0.0.0", "ip of the interface for the server to listen on")
	period = flag.Duration("sync-duration", time.Second, "duration after which append only file will be flushed to disk")
)

type Config struct {
	Host   string
	Port   int
	Period time.Duration
}

func NewConfig() *Config {
	flag.Parse()
	return &Config{
		Port:   *port,
		Host:   *host,
		Period: *period,
	}
}
