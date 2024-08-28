package client

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Client struct {
	client     http.Client
	logger     zerolog.Logger
	serverAddr string
	name       string
	closeChan  chan any
	closedChan chan any
}

func (c *Client) run(interval time.Duration) {
	defer close(c.closedChan)
	defer c.logger.Info().Msg("client stopped")
	c.logger.Info().Msg("client started")
	var lastRun time.Time
	for {
		var (
			nextRun = lastRun.Add(interval)
			now     = time.Now()
		)
		if nextRun.Before(now) {
			if err := c.update(); err != nil {
				c.logger.Error().Msg(err.Error())
			}
			lastRun = now
			nextRun = now.Add(interval)
		}
		t := time.NewTimer(nextRun.Sub(now))
		select {
		case <-t.C:
		case <-c.closeChan:
			t.Stop()
			return
		}
	}
}

func New(cfg *Config) (*Client, error) {

	// Load the certificate and private key
	cert, err := tls.LoadX509KeyPair(
		cfg.CertFilename,
		cfg.KeyFilename,
	)
	if err != nil {
		return nil, err
	}

	// Load the CA certificate
	caCert, err := os.ReadFile(cfg.CACertFilename)
	if err != nil {
		return nil, err
	}

	// Create a certificate pool and add the CA certificate
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	// Create the client instance
	c := &Client{
		client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
					RootCAs:      certPool,
				},
			},
		},
		logger:     log.With().Str("package", "client").Logger(),
		serverAddr: cfg.ServerAddr,
		name:       cfg.Name,
		closeChan:  make(chan any),
		closedChan: make(chan any),
	}

	go c.run(cfg.Interval)

	return c, nil
}

func (c *Client) Close() {
	close(c.closeChan)
	<-c.closedChan
}
