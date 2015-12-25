package grpcutil

import (
	"net"
	"os"

	"google.golang.org/grpc"
)

type localServer struct {
	server *grpc.Server
	path   string
}

func (s *localServer) Server() *grpc.Server {
	return s.server
}

func (s *localServer) Serve() (retErr error) {
	listener, err := net.Listen("unix", s.path)
	if err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(s.path); err != nil && retErr == nil {
			retErr = err
		}
	}()
	return s.server.Serve(listener)
}

func (s *localServer) Dial() (*grpc.ClientConn, error) {
	return grpc.Dial("unix://" + s.path)
}
