package main

import (
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v3"
)

func readConfigFile(filename string, v any) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(v)
}

func waitForInterrupt() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
