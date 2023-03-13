package main

import (
	"os"
	"strconv"

	"clientservice/service"

	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().
		Str("server", "client").Logger()

	// Get HTTP port from environment variable
	portString := os.Getenv("HTTP_PORT")
	if portString == "" {
		logger.Fatal().Msg("Please specify HTTP port in HTTP_PORT environment variable")
	}

	port, err := strconv.ParseInt(portString, 10, 64)
	if err != nil {
		logger.Fatal().Err(err).Msg("error parsing HTTP_PORT")
	}

	// Connect to Ports gRPC server
	portsSrv, portsConn, err := service.NewPortsService()
	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to Ports gRPC server")
	}

	err = service.StartHTTPServer(&logger, int(port), portsSrv, portsConn)
	if err != nil {
		logger.Fatal().Err(err).Msg("error starting http client server")
	}
}
