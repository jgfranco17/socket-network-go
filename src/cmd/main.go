package main

import (
	"flag"
	"time"

	"src/logging"
	"src/server"

	log "github.com/sirupsen/logrus"
)

var (
	port     *string = flag.String("port", "8080", "Port to listen on")
	lifespan *int    = flag.Int("lifespan", 60, "Port to listen on")
)

func init() {
	logging.ConfigureLogger("DEBUG")
}

func main() {
	flag.Parse()
	duration := time.Duration(*lifespan) * time.Second
	server := server.NewServer(*port, duration)
	err := server.Run()
	if err != nil {
		log.Fatalf("Server failed unexpectedly: %v", err)
	}
}
