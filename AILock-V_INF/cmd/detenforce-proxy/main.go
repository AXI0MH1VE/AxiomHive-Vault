package main

import (
	"flag"
	"log"
	"net"
	"net/http"
)

func main() {
	port := flag.String("port", "8443", "The server port")
	configPath := flag.String("config", "/app/config/proxy_core.yaml", "Config file path")
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// TODO: Implement DetEnforceProxy with TLS termination, validation, rate limiting
	// Forward requests to AuthService via gRPC

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Placeholder handler
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("DetEnforceProxy running"))
		}),
	}

	log.Printf("Starting DetEnforceProxy on port %s with config %s", *port, *configPath)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
