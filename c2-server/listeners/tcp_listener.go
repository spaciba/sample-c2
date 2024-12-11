package listeners

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TCPServer struct {
	address    string
	listener   net.Listener
	shutdownCh chan struct{}
	wg         *sync.WaitGroup
}

func NewTCPServer(address string, wg *sync.WaitGroup) *TCPServer {
	return &TCPServer{
		address:    address,
		shutdownCh: make(chan struct{}),
		wg:         wg,
	}
}

func (s *TCPServer) Start() {
	defer s.wg.Done()

	// Start listening on the given address
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("Failed to start server on %s: %v", s.address, err)
		return
	}
	s.listener = listener
	log.Printf("Server started on %s", s.address)

	// Accept incoming connections in a loop
	for {
		select {
		case <-s.shutdownCh:
			log.Printf("Shutting down server on %s", s.address)
			s.listener.Close()
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting connection on %s: %v", s.address, err)
				continue
			}
			go handleConnection(conn)
		}
	}
}

func (s *TCPServer) Stop() {
	close(s.shutdownCh)
}

// handleConnection handles incoming client connections
func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("New connection from %s", conn.RemoteAddr())

	// Simulate some work with the connection
	time.Sleep(2 * time.Second)
	_, err := conn.Write([]byte("Hello from server!"))
	if err != nil {
		log.Printf("Failed to send data: %v", err)
	}
	log.Printf("Closed connection from %s", conn.RemoteAddr())
}

type TCPListener struct {
	listener_address string
	listener_state   ListenerState
	tcp_server       *TCPServer
	listener_id string
}

func (tcp *TCPListener) Init(listener_address string) (listenerID string, message string, error_code int32) {
	var wg sync.WaitGroup

	server := NewTCPServer(listener_address, &wg)

	id := uuid.New()

	tcp = &TCPListener{listener_address: listener_address, listener_state: STARTED, tcp_server: server}

	return id.String(), "TCP Listener Created!", 0
}
