package main

import (
	"log"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	server := NewWeatherServer(config)

	log.Fatal(server.Start())
}
