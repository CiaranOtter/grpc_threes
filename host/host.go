package main

import (
	"context"
	"fmt"
	"log"
	"net"

	gm "github.com/CiaranOtter/grpc_threes.git/gameservice"
	"google.golang.org/grpc"
)

type GameService struct {
	gm.UnimplementedGameServiceServer
	// players
}

func (s *GameService) MakeMove(ctx context.Context, req *gm.Move) (*gm.Response, error) {

	return &gm.Response{}, nil
}

func (s *GameService) GetCommand(ctx context.Context, req *gm.Empty) (*gm.Command, error) {
	return &gm.Command{}, nil
}

func (s *GameService) RegisterPlayer(ctx context.Context, req *gm.Player) (*gm.Player, error) {
	fmt.Printf("New player with the name %s has entered the game\n", req.GetName())
	return &gm.Player{}, nil
}

func StartService(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	service := &GameService{}
	grpcServer := grpc.NewServer(opts...)
	gm.RegisterGameServiceServer(grpcServer, service)

	grpcServer.Serve(lis)
}

func main() {
	StartService(5000)
}
