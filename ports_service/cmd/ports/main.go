package main

import (
	"portsservice/service"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	log := logger.WithFields(logrus.Fields{
		"Server": "Ports GRPC Server",
	})
	service.StartGRPCServer(log)
}
