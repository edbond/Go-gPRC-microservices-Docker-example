package main

import (
	"os"
	"strconv"

	"clientservice"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	log := logger.WithFields(logrus.Fields{
		"Server": "Client",
	})

	// Get HTTP port from environment variable
	portString := os.Getenv("HTTP_PORT")
	if portString == "" {
		log.Panicln("Please specify HTTP port in HTTP_PORT environment variable")
	}

	port, err := strconv.ParseInt(portString, 10, 64)
	if err != nil {
		log.Panicf("error parsing HTTP_PORT: %s", err)
	}

	// Connect to Ports gRPC server
	portsSrv, portsConn, err := clientservice.NewPortsService()
	if err != nil {
		log.Panicln("error connecting to Ports gRPC server", err)
	}

	err = clientservice.StartHTTPServer(log, int(port), portsSrv, portsConn)
	if err != nil {
		log.Panicln("error starting http client server", err)
	}
}
