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
	go func() {
		time.Sleep(sig.Duration)
		close(sig.ShutdownChannel)
	}()
}

type Server struct {
	port           string
	shutdownSignal *ShutdownSignal
	startTime      time.Time
	lifespan       time.Duration
}

// NewServer creates a new UDP server
func NewServer(port string, lifespan time.Duration) *Server {
	return &Server{
		port:           port,
		shutdownSignal: NewShutdownSignal(lifespan),
		startTime:      time.Now(),
		lifespan:       lifespan,
	}
}

func (s *Server) Address() string {
	return fmt.Sprintf("localhost:%s", s.port)
}

func (s *Server) CheckShutdown() chan struct{} {
	return s.shutdownSignal.ShutdownChannel
}

// Start initializes the server and starts listening for UDP messages
func (s *Server) Run() error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("Error resolving address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
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
			n, clientAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				netErr, ok := err.(net.Error)
				if ok && netErr.Timeout() {
					// Timeout expired, go back to listening (ignore timeout errors)
					continue
				}
				PrintError("Error reading from client: %v", err)
				continue
			}

			message := strings.TrimSpace(string(buffer[:n]))
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
	if _, err := conn.WriteToUDP([]byte(response), senderAddr); err != nil {
		log.Infof("Error sending message: %v", err)
	}
}
