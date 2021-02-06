package service

import (
	context "context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"portsservice/internal"
	"portsservice/repository"
	"syscall"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	ports "ports.services.com/ports"
)

type portsService struct {
	ports.UnimplementedPortsServiceServer
	repository repository.PortsRepository
}

func (srv *portsService) List(req *ports.ListRequest, listResponse ports.PortsService_ListServer) error {
	allPorts, err := srv.repository.AllPorts()
	if err != nil {
		return err
	}

	for _, p := range allPorts {
		err = listResponse.Send(ports.PortToProto(&p))
		if err != nil {
			return err
		}
	}
	return nil
}

func (srv *portsService) Upsert(ctx context.Context, port *ports.PortProto) (*ports.UpsertResponse, error) {
	portStruct := ports.ProtoToPort(port)

	err := srv.repository.Upsert(*portStruct)
	return &ports.UpsertResponse{}, err
}

// StartGRPCServer starts Ports GRPC server
func StartGRPCServer(log *logrus.Entry) error {
	grpcServer := grpc.NewServer()

	grpcService := portsService{}
	repository := internal.MemoryRepository{}
	err := repository.Init()
	if err != nil {
		return fmt.Errorf("repository initialization error: %w", err)
	}
	grpcService.repository = &repository

	ports.RegisterPortsServiceServer(grpcServer, &grpcService)

	grpcPort := os.Getenv("PORTS_GRPC_PORT")
	if grpcPort == "" {
		log.Panic("Please specify tcp port for grpc server in PORTS_GRPC_PORT envrironment variable")
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
		grpcService.repository.Close()
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
