package config

import "flag"

var (
	port = flag.Int("port", 6379, "port for redis-lite to listen on")
	host = flag.String("ip", "0.0.0.0", "ip of the interface for the server to listen on")
)

type Config struct {
	Host string
	Port int
}

func NewConfig() *Config {
	flag.Parse()
	return &Config{
		Port: *port,
		Host: *host,
	}
}
