package main

import (
	"os"
	"os/signal"
	"syscall"
)

func runSignal() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	return nil
}
