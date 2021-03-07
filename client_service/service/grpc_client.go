package service

import (
	"errors"
	"fmt"
	"os"

	"clientservice/ports"

	"google.golang.org/grpc"
)

// NewPortsService connects to Ports gRPC server
// and returns service, connection and error
func NewPortsService() (ports.PortsServiceClient, *grpc.ClientConn, error) {
	portsServiceAddress := os.Getenv("PORTS_ADDRESS")
	if portsServiceAddress == "" {
		return nil, nil, errors.New("please specify address of ports service in PORTS_ADDRESS environment variable")
	}

	portsConn, err := grpc.Dial(portsServiceAddress, grpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("Can't connect to Ports service using address %s: %s", portsServiceAddress, err)
	}

	portsSrv := ports.NewPortsServiceClient(portsConn)
	return portsSrv, portsConn, nil
}
