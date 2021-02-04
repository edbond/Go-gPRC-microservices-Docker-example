package clientservice

import "google.golang.org/grpc"

func newPortsService(target string) (*grpc.ClientConn, error) {
	return grpc.Dial(target, grpc.WithInsecure())
}
