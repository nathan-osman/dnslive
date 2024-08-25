package server

type Config struct {
	CACertFilename string `yaml:"ca_cert_filename"`
	CertFilename   string `yaml:"cert_filename"`
	KeyFilename    string `yaml:"key_filename"`
	HttpServerAddr string `yaml:"http_server_addr"`
	DnsServerAddr  string `yaml:"dns_server_addr"`
	PersistentFile string `yaml:"persistent_file"`
}
