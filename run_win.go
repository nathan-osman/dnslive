//go:build windows

package main

import (
	"golang.org/x/sys/windows/svc"
)

const (
	serviceName = "dnslive"
)

type handlerFunc func([]string, <-chan svc.ChangeRequest, chan<- svc.Status) (bool, uint32)

func (f handlerFunc) Execute(
	args []string,
	chChan <-chan svc.ChangeRequest,
	stChan chan<- svc.Status,
) (bool, uint32) {
	return f(args, chChan, stChan)
}

func runService(
	args []string,
	chChan <-chan svc.ChangeRequest,
	stChan chan<- svc.Status,
) (bool, uint32) {

	// Indicate that the service has been started
	stChan <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop | svc.AcceptShutdown,
	}

	// Respond to service requests
	for c := range chChan {
		switch c.Cmd {
		case svc.Interrogate:
			stChan <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			stChan <- svc.Status{
				State: svc.StopPending,
			}
			return false, 0
		}
	}

	// This line should never be reached
	return false, 0
}

func run() error {

	// See if application is being run interactively
	if i, err := svc.IsWindowsService(); err != nil {
		return err
	} else if !i {
		return runSignal()
	}

	return svc.Run(serviceName, handlerFunc(runService))
}
