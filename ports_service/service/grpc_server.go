package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"portsservice/ports"
	"portsservice/repository"
	"portsservice/storage"
	"syscall"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type portsService struct {
	ports.UnimplementedPortsServiceServer
	repository repository.PortsRepository
}

func (srv *portsService) List(_ *ports.ListRequest, listResponse ports.PortsService_ListServer) error {
	allPorts, err := srv.repository.AllPorts()
	if err != nil {
		return err
	}

	for i := range allPorts {
		t := allPorts[i].ToTransport()
		err = listResponse.Send(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (srv *portsService) Upsert(_ context.Context, port *ports.PortTransport) (*ports.UpsertResponse, error) {
	err := srv.repository.Upsert(port.ToValue())
	return &ports.UpsertResponse{}, err
}

// StartGRPCServer starts Ports GRPC server
func StartGRPCServer(log *logrus.Entry) error {
	grpcServer := grpc.NewServer()

	grpcService := portsService{}
	memoryRepository := storage.MemoryRepository{}
	err := memoryRepository.Init()
	if err != nil {
		return fmt.Errorf("memoryRepository initialization error: %w", err)
	}
	grpcService.repository = &memoryRepository

	ports.RegisterPortsServiceServer(grpcServer, &grpcService)

	grpcPort := os.Getenv("PORTS_GRPC_PORT")
	if grpcPort == "" {
		log.Panic("Please specify tcp port for grpc server in PORTS_GRPC_PORT environment variable")
	}

	con, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Panic(err)
	}

	log.Infof("Starting Ports GRPC Server on port %s", grpcPort)

	// Listen for signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	shutdownComplete := make(chan struct{}, 1)
	go func() {
		<-c

		// Shutdown GRPC server!
		log.Infoln("🔻 Shutdown server")
		err := grpcService.repository.Close()
		if err != nil {
			log.Errorf("error closing repository %s", err)
		}
		grpcServer.GracefulStop()
		log.Infoln("Stopped successfully")
		shutdownComplete <- struct{}{}
	}()

	err = grpcServer.Serve(con)
	if err != nil {
		log.Panic(err)
	}
	<-shutdownComplete

	return nil
}
