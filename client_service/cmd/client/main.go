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

	portString := os.Getenv("HTTP_PORT")
	if portString == "" {
		log.Panicln("Please specify HTTP port in HTTP_PORT environment variable")
	}

	port, err := strconv.ParseInt(portString, 10, 32)
	if err != nil {
		log.Panicf("error parsing HTTP_PORT: %s", err)
	}

	err = clientservice.StartHTTPServer(log, int(port))
	if err != nil {
		log.Errorln("error starting http client server", err)
	}
}
