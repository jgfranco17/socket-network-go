package main

import (
	"flag"

	"src/logging"
	"src/server"

	log "github.com/sirupsen/logrus"
)

var (
	port     *string = flag.String("port", "8080", "Port to listen on")
	lifespan *int    = flag.Int("lifespan", 0, "Port to listen on")
)

func init() {
	logging.ConfigureLogger("DEBUG")
}

func main() {
	flag.Parse()
	server, err := server.CreateNewServerUDP(*port, *lifespan)
	if err != nil {
		log.Fatalf("Server failed to initialize: %v", err)
	}
	err = server.Start()
	if err != nil {
		log.Fatalf("Server error during run: %v", err)
	}
}
