package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"clientservice/ports"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// StartHTTPServer loads json and starts HTTP server
func StartHTTPServer(logger *zerolog.Logger, port int, portsSrv ports.PortsServiceClient, portsConn *grpc.ClientConn) error {
	var err error

	filename := os.Getenv("PORTS_JSON")
	if filename == "" {
		logger.Fatal().Msg("Please set PORTS_JSON env variable to path to ports json file")
	}

	err = loadJSON(logger, filename, portsSrv)
	if err != nil {
		panic(fmt.Errorf("error loading json file: %s", err.Error()))
	}

	router := http.NewServeMux()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
		// timeouts for bad http clients
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	// HTTP Handler
	// GET / returns all ports from ports service
	router.HandleFunc("/ports", func(rw http.ResponseWriter, r *http.Request) {
		getPorts(logger, portsSrv, rw, r)
	})

	// Listen for signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Info().Msg("🔻 Shutdown server")
		err := srv.Shutdown(context.Background())
		if err != nil {
			logger.Err(err).Msg("server shutdown error")
		}

		// Shutdown connection to Ports Service - Bye!
		err = portsConn.Close()
		if err != nil {
			logger.Err(err).Msg("error closing connection to ports service")
		}
	}()

	logger.Info().Msg("🆙 Starting server at port")

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("error starting client server: %w", err)
	}

	return nil
}

func loadJSON(logger *zerolog.Logger, filename string, portsSrv ports.PortsServiceClient) error {
	var err error

	totalPorts := 0
	failedPorts := 0

	ctx := context.Background()

	callback := func(port *ports.Port) {
		_, err = portsSrv.Upsert(ctx, port)
		if err != nil {
			logger.Err(err).Msg("error sending Port to Ports service")
			failedPorts++
		}
		totalPorts++
	}

	err = ports.LoadFromJSON(logger, filename, callback)
	if err != nil {
		return err
	}

	logger.Info().Msgf("Total ports loaded: %d, failed: %d", totalPorts, failedPorts)

	return nil
}

func getPorts(logger *zerolog.Logger, portsSrv ports.PortsServiceClient, w http.ResponseWriter, _ *http.Request) {
	var err error

	ctx := context.Background()
	allPortsClient, err := portsSrv.List(ctx, &ports.ListRequest{})
	if err != nil {
		httpError(logger, err, "error from Ports service", w)
		return
	}

	var allPorts []*ports.Port
	for {
		port, err := allPortsClient.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			httpError(logger, err, "error reading Port from service", w)
			return
		}

		allPorts = append(allPorts, port)
	}

	portsJSON, err := json.Marshal(allPorts)
	if err != nil {
		httpError(logger, err, "error serializing JSON", w)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(portsJSON); err != nil {
		logger.Warn().Err(err).Msg("error writing output json")
	}
}

func httpError(logger *zerolog.Logger, err error, msg string, w http.ResponseWriter) {
	errMessage := fmt.Sprintf("%s: %s", msg, err)

	w.WriteHeader(http.StatusBadGateway)
	if _, err := w.Write([]byte(errMessage)); err != nil {
		logger.Warn().Err(err).Msg("error writing output json")
	}
}
