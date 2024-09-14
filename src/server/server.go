package server

import (
	"fmt"
	"net"
	"strings"
	"time"

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

// StartShutdownTimer starts the shutdown timer
func (sig *ShutdownSignal) StartShutdownTimer() {
	if sig.Duration == 0 {
		log.Warnf("No shutdown timer set, server will run indefinitely.")
		return
	}
	go func() {
		time.Sleep(sig.Duration)
		close(sig.ShutdownChannel)
	}()
}

type Server struct {
	address        net.Addr
	port           string
	shutdownSignal *ShutdownSignal
	startTime      time.Time
	lifespan       time.Duration
}

// Creates a new UDP server
func CreateNewServerUDP(port string, lifespan time.Duration) (*Server, error) {
	if lifespan < 0 {
		return nil, fmt.Errorf("Lifetime cannot be negative.")
	}
	addr, err := net.ResolveUDPAddr(udpNetworkType, fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, fmt.Errorf("Error resolving address: %w", err)
	}
	return &Server{
		address:        addr,
		port:           port,
		shutdownSignal: NewShutdownSignal(lifespan),
		startTime:      time.Now(),
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
func (s *Server) Run() error {
	addr, err := net.ResolveUDPAddr(udpNetworkType, fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("Error resolving address: %w", err)
	}

	conn, err := net.ListenUDP(udpNetworkType, addr)
	if err != nil {
		return fmt.Errorf("Error starting server: %w", err)
	}
	defer conn.Close()

	s.shutdownSignal.StartShutdownTimer()
	buffer := make([]byte, maxByteSize)
	log.Infof("Server listening on %s for %vs", s.Address(), s.lifespan.Seconds())

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
				PrintError("Error reading from client: %v", err)
				continue
			}

			message := strings.TrimSpace(string(buffer[:idx]))
			if message == "" {
				continue
			}

			PrintStandardMessage(clientAddr.String(), message)
			s.BroadcastMessage(message, conn, clientAddr)
		}
	}
}

// broadcastMessage sends the message back to the client who sent it
func (s *Server) BroadcastMessage(message string, conn *net.UDPConn, senderAddr *net.UDPAddr) {
	response := fmt.Sprintf("[%s]: %s", senderAddr.String(), message)

	// Send the message back to the client who sent it
	_, err := conn.WriteToUDP([]byte(response), senderAddr)
	if err != nil {
		log.Infof("Error sending response message: %v", err)
	}
}
