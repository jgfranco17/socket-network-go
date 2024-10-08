package server

import (
	"fmt"
	"net"
	"strings"
	"time"

	"src/outputs"

	log "github.com/sirupsen/logrus"
)

type ShutdownSignal struct {
	ShutdownChannel chan struct{}
	Duration        time.Duration
}

func NewShutdownSignal(duration time.Duration) *ShutdownSignal {
	return &ShutdownSignal{
		ShutdownChannel: make(chan struct{}),
		Duration:        duration,
	}
}

type Server struct {
	address        *net.UDPAddr
	port           string
	shutdownSignal *ShutdownSignal
	startTime      time.Time
	lifespan       time.Duration
}

func (s *Server) startShutdownTimer() {
	go func() {
		time.Sleep(s.lifespan)
		close(s.shutdownSignal.ShutdownChannel)
	}()
}

// Creates a new UDP server
func CreateNewServerUDP(port string, lifespanMinutes int) (*Server, error) {
	if port == "" {
		return nil, fmt.Errorf("No port specified")
	}
	if lifespanMinutes < 0 {
		return nil, fmt.Errorf("Lifetime cannot be negative")
	}
	addr, err := net.ResolveUDPAddr(udpNetworkType, fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, fmt.Errorf("Error resolving address: %w", err)
	}
	lifespan := time.Duration(lifespanMinutes) * time.Minute
	return &Server{
		address:        addr,
		port:           port,
		shutdownSignal: NewShutdownSignal(lifespan),
		startTime:      time.Now().Local(),
		lifespan:       lifespan,
	}, nil
}

func (s *Server) Address() string {
	return s.address.String()
}

func (s *Server) CheckShutdown() chan struct{} {
	return s.shutdownSignal.ShutdownChannel
}

// Start initializes the server and starts listening for UDP messages
func (s *Server) Start() error {
	conn, err := net.ListenUDP(udpNetworkType, s.address)
	if err != nil {
		return fmt.Errorf("Error starting server: %w", err)
	}
	defer conn.Close()

	buffer := make([]byte, maxByteSize)
	if s.shutdownSignal.Duration == 0 {
		log.Warnf("No shutdown timer set, server will run indefinitely.")
	} else {
		s.startShutdownTimer()
		log.Infof("Server listening on %s for %vm", s.Address(), s.lifespan.Minutes())
	}

	for {
		select {
		case <-s.CheckShutdown():
			log.Infof("Server shutting down...")
			return nil
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			idx, clientAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				netErr, ok := err.(net.Error)
				if ok && netErr.Timeout() {
					// Timeout expired, go back to listening (ignore timeout errors)
					continue
				}
				outputs.PrintError("Error reading from client: %v", err)
				continue
			}

			message := strings.TrimSpace(string(buffer[:idx]))
			if message == "" {
				continue
			}
			s.BroadcastMessage(message, conn, clientAddr)
		}
	}
}

// broadcastMessage sends the message back to the client who sent it
func (s *Server) BroadcastMessage(message string, conn *net.UDPConn, addr *net.UDPAddr) {
	outputs.PrintStandardMessage(addr.String(), message)
	response := fmt.Sprintf("[%s]: %s", addr.String(), message)

	// Send the message back to the client who sent it
	_, err := conn.WriteToUDP([]byte(response), addr)
	if err != nil {
		log.Infof("Error sending response message: %v", err)
	}
}
