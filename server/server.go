package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type entry struct {
	LastUpdate time.Time `json:"last_update"`
	Ipv4       string    `json:"ipv4"`
	Ipv6       string    `json:"ipv6"`
}

type Server struct {
	mutex      sync.RWMutex
	httpServer http.Server
	dnsServer  dns.Server
	logger     zerolog.Logger
	filename   string
	entries    map[string]entry
}

func New(cfg *Config) (*Server, error) {

	// Load the CA certificate
	caCert, err := os.ReadFile(cfg.CACertFilename)
	if err != nil {
		return nil, err
	}

	// Create a certificate pool and add the CA certificate
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	// Create the server instance
	var (
		r http.ServeMux
		h dns.ServeMux
		s = &Server{
			httpServer: http.Server{
				Addr:    cfg.HttpServerAddr,
				Handler: &r,
				TLSConfig: &tls.Config{
					ClientCAs:  certPool,
					ClientAuth: tls.RequireAndVerifyClientCert,
				},
			},
			dnsServer: dns.Server{
				Addr:      cfg.DnsServerAddr,
				Net:       "udp",
				Handler:   &h,
				ReusePort: true,
			},
			logger:   log.With().Str("package", "server").Logger(),
			filename: cfg.PersistentFile,
			entries:  make(map[string]entry),
		}
	)

	// Load the existing entries
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Handle updates from clients
	r.HandleFunc("/update", s.update)

	// Handle incoming DNS requests
	h.HandleFunc(".", s.respond)

	// Listen for DNS requests
	go func() {
		defer s.logger.Info().Msg("DNS server stopped")
		s.logger.Info().Msg("DNS server starting...")
		if err := s.dnsServer.ListenAndServe(); err != nil {
			s.logger.Error().Msg(err.Error())
		}
	}()

	// Listen for HTTP connections
	go func() {
		defer s.logger.Info().Msg("HTTP server stopped")
		s.logger.Info().Msg("HTTP server starting...")
		if err := s.httpServer.ListenAndServeTLS(
			cfg.CertFilename,
			cfg.KeyFilename,
		); !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error().Msg(err.Error())
		}
	}()

	return s, nil
}

func (s *Server) Close() {
	s.httpServer.Shutdown(context.Background())
	s.dnsServer.Shutdown()
}
