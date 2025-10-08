package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"os"
	"strings"
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
	mutex        sync.RWMutex
	httpServer   http.Server
	dnsTCPServer dns.Server
	dnsUDPServer dns.Server
	logger       zerolog.Logger
	zone         string
	nameserver   string
	mailbox      string
	filename     string
	entries      map[string]entry
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
			dnsTCPServer: dns.Server{
				Addr:    cfg.DnsServerAddr,
				Net:     "tcp",
				Handler: &h,
			},
			dnsUDPServer: dns.Server{
				Addr:      cfg.DnsServerAddr,
				Net:       "udp",
				Handler:   &h,
				ReuseAddr: true,
				ReusePort: true,
			},
			logger:     log.With().Str("package", "server").Logger(),
			zone:       cfg.Zone,
			nameserver: cfg.Nameserver,
			mailbox:    cfg.Mailbox,
			filename:   cfg.PersistentFile,
			entries:    make(map[string]entry),
		}
	)

	// Ensure names end with a "."
	if !strings.HasSuffix(s.zone, ".") {
		s.zone += "."
	}
	if !strings.HasSuffix(s.nameserver, ".") {
		s.nameserver += "."
	}
	if !strings.HasSuffix(s.mailbox, ".") {
		s.mailbox += "."
	}

	// Load the existing entries
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Handle updates from clients
	r.HandleFunc("/update", s.update)

	// Handle incoming DNS requests
	h.HandleFunc(".", s.respond)

	// Listen for TCP DNS requests
	go func() {
		defer s.logger.Info().Msg("TCP DNS server stopped")
		s.logger.Info().Msg("TCP DNS server starting...")
		if err := s.dnsTCPServer.ListenAndServe(); err != nil {
			s.logger.Error().Msg(err.Error())
		}
	}()

	// Listen for UDP DNS requests
	go func() {
		defer s.logger.Info().Msg("UDP DNS server stopped")
		s.logger.Info().Msg("UDP DNS server starting...")
		if err := s.dnsUDPServer.ListenAndServe(); err != nil {
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
	s.dnsTCPServer.Shutdown()
	s.dnsUDPServer.Shutdown()
}
