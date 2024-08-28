package client

import (
	"time"
)

type Config struct {
	CACertFilename string        `yaml:"ca_cert_filename"`
	CertFilename   string        `yaml:"cert_filename"`
	KeyFilename    string        `yaml:"key_filename"`
	ServerAddr     string        `yaml:"server_addr"`
	Interval       time.Duration `yaml:"interval"`
	Name           string        `yaml:"name"`
}
