package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/spaciba/sample_c2/api"
	"github.com/spaciba/sample_c2/c2-server/listeners"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedControllerServer
	State *ServerState
}

// SayHello implements helloworld.GreeterServer
func (s *server) CreateListener(_ context.Context, in *pb.CreateListenerRequest) (*pb.CreateListenerReply, error) {
	log.Printf("Creating %v listener on port %v", in.GetListenerType(), in.GetListenerAddress())

	var listener listeners.Listener
	var listener_id string
	var message string
	var error_code int32

	if in.ListenerType == "tcp" {
		listener = &listeners.TCPListener{}
		listener_id, message, error_code = listener.Init(in.GetListenerAddress())
		s.State.listeners[listener_id] = &listener
	}

	fmt.Printf("State: %v", s.State.listeners)
	
	return &pb.CreateListenerReply{Message: message, ListenerId: listener_id, ErrorCode: error_code}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	state := NewServerState()
	pb.RegisterControllerServer(s, &server{State: state})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}