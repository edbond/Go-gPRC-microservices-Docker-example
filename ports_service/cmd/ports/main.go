package main

import (
	"portsservice/service"

	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(zerolog.NewConsoleWriter())

	if err := service.StartGRPCServer(&logger); err != nil {
		logger.Panic().Err(err).Msg("error running grpc server")
	}
}
