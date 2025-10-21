package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "50052", "The server port")
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// TODO: Generate proto code and implement RBAC service
	s := grpc.NewServer()

	log.Printf("Starting RBACService on port %s", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
