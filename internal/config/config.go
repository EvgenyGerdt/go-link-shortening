package config

import "flag"

type Config struct {
	ServerAddress string
	BaseURL       string
}

func NewConfig() Config {
	serverAddress := flag.String("a", ":8080", "Server address")
	baseURL := flag.String("b", "http://localhost:8080", "Base URL")

	flag.Parse()

	return Config{
		ServerAddress: *serverAddress,
		BaseURL:       *baseURL,
	}
}
