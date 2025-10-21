package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		ListenAddress string `yaml:"listen_address"`
		Port          string `yaml:"port"`
		TLS           struct {
			Enabled  bool   `yaml:"enabled"`
			CertFile string `yaml:"cert_file"`
			KeyFile  string `yaml:"key_file"`
		} `yaml:"tls"`
	} `yaml:"server"`
	Upstreams []struct {
		Name    string `yaml:"name"`
		Address string `yaml:"address"`
	} `yaml:"upstreams"`
	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

func main() {
	configPath := flag.String("config", "/app/config/proxy_core.yaml", "Config file path")
	flag.Parse()

	// Load configuration
	configData, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Setup logging
	if cfg.Logging.Output != "" {
		logFile, err := os.OpenFile(cfg.Logging.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file %s: %v", cfg.Logging.Output, err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	if len(cfg.Upstreams) == 0 {
		log.Fatalf("No upstreams defined in config file")
	}
	upstreamURL, err := url.Parse(fmt.Sprintf("http://%s", cfg.Upstreams[0].Address))
	if err != nil {
		log.Fatalf("Failed to parse upstream URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(upstreamURL)

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.Server.ListenAddress, cfg.Server.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("INFO: Proxying request: %s %s to %s", r.Method, r.URL.Path, upstreamURL)
			proxy.ServeHTTP(w, r)
		}),
	}

	log.Printf("INFO: Starting DetEnforceProxy on %s with config %s", server.Addr, *configPath)

	if cfg.Server.TLS.Enabled {
		log.Printf("TLS enabled: using cert %s and key %s", cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile)
		if err := server.ListenAndServeTLS(cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile); err != nil {
			log.Fatalf("Failed to serve with TLS: %v", err)
		}
	} else {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}
}

