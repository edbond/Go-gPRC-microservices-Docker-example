package clientservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
	"ports.services.com/ports"
)

// StartHTTPServer loads json and starts HTTP server
func StartHTTPServer(log *logrus.Entry, port int) error {
	var err error

	filename := os.Getenv("PORTS_JSON")
	if filename == "" {
		log.Panic("Please set PORTS_JSON env variable to path to ports json file")
	}

	portsServiceAddress := os.Getenv("PORTS_ADDRESS")
	if portsServiceAddress == "" {
		log.Panic("Please specify address of ports service in PORTS_ADDRESS environment variable")
	}

	portsConn, err := newPortsService(portsServiceAddress)
	if err != nil {
		log.Panicf("Can't connect to Ports service using address %s: %s", portsServiceAddress, err)
	}

	portsSrv := ports.NewPortsServiceClient(portsConn)

	err = loadJSON(log, filename, portsSrv)
	if err != nil {
		panic(fmt.Errorf("error loading json file: %s", err.Error()))
	}

	router := http.NewServeMux()

	srv := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: router,
		// timeouts for bad http clients
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	// HTTP Hanlder
	// GET / returns all ports from ports service
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		getPorts(portsSrv, rw, r)
	})

	// Listen for signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		log.Infoln("🔻 Shutdown server")
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Errorf("server shutdown error: %s", err)
		}

		// Shutdown connection to Ports Service - Bye!
		err = portsConn.Close()
		if err != nil {
			log.Errorf("error closing connection to ports service: %s", err)
		}
	}()

	log.Infof("🆙 Starting server at port %d", port)

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("error starting client server: %w", err)
	}

	return nil
}

func loadJSON(log *logrus.Entry, filename string, portsSrv ports.PortsServiceClient) error {
	var err error

	totalPorts := 0
	failedPorts := 0

	ctx := context.Background()

	callback := func(port *ports.Port) {
		// log.Infof("Port: %v\n", port)

		_, err = portsSrv.Upsert(ctx, ports.PortToProto(port))
		if err != nil {
			log.Errorf("error sending Port to Ports service: %s", err)
			failedPorts++
		}
		totalPorts++
	}

	err = ports.LoadFromJSON(log, filename, callback)
	if err != nil {
		return err
	}

	log.Infof("Total ports loaded: %d, failed: %d", totalPorts, failedPorts)

	return nil
}

func getPorts(portsSrv ports.PortsServiceClient, w http.ResponseWriter, req *http.Request) {

	ctx := context.Background()
	allPortsClient, err := portsSrv.List(ctx, &ports.ListRequest{})
	if err != nil {
		httpError(err, "error from Ports service", w)
		return
	}

	allPorts := []ports.Port{}
	for {
		port, err := allPortsClient.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			httpError(err, "error reading Port from service", w)
			return
		}

		allPorts = append(allPorts, *ports.ProtoToPort(port))
	}

	portsJSON, err := json.Marshal(allPorts)
	if err != nil {
		httpError(err, "error serializing JSON", w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(portsJSON)
}

func httpError(err error, msg string, w http.ResponseWriter) {
	errMessage := fmt.Sprintf("%s: %s", msg, err)

	w.WriteHeader(http.StatusBadGateway)
	w.Write([]byte(errMessage))
}
