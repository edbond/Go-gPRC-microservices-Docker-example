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

	if err := service.StartGRPCServer(log); err != nil {
		log.Panic(err)
	}
}
