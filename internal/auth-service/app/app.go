package app

import (
	"net"

	"github.com/Dorrrke/note-tracker/internal/auth-service/server"
	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
}

func NewApp() *App {
	grpcServe := configureGrpcServer()
	return &App{
		gRPCServer: grpcServe,
	}
}

func (a *App) StartApp() error {
	return runServer(a.gRPCServer)
}

func configureGrpcServer() *grpc.Server {
	grpcServe := grpc.NewServer()
	server.RegisterGrpcServer(grpcServe)
	return grpcServe
}

func runServer(grpServe *grpc.Server) error {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		return err
	}
	if err := grpServe.Serve(listener); err != nil {
		return err
	}
	return nil
}
