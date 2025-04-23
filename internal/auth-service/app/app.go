package app

import (
	"fmt"
	"net"

	"github.com/Dorrrke/note-tracker/internal/auth-service/config"
	"github.com/Dorrrke/note-tracker/internal/auth-service/server"
	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
	cfg        *config.Config
}

func NewApp(cfg *config.Config, repo server.Repository) *App {
	grpcServe := configureGrpcServer(repo)
	return &App{
		gRPCServer: grpcServe,
		cfg:        cfg,
	}
}

func (a *App) StartApp() error {
	return runServer(a.gRPCServer, a.cfg.Host, a.cfg.Port)
}

func configureGrpcServer(repo server.Repository) *grpc.Server {
	grpcServe := grpc.NewServer()
	server.RegisterGrpcServer(grpcServe, repo)
	return grpcServe
}

func runServer(grpServe *grpc.Server, host string, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	if err := grpServe.Serve(listener); err != nil {
		return err
	}
	return nil
}
