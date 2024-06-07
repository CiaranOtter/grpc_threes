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
	player []string
}

func (s *GameService) MakeMove(ctx context.Context, req *gm.Move) (*gm.Response, error) {
	fmt.Printf("Player with the colour %d made move %d, %d\n", req.Colour, req.Row, req.Col)
	return &gm.Response{
		Success: true,
	}, nil
}

func (s *GameService) GetCommand(ctx context.Context, req *gm.Empty) (*gm.Command, error) {
	fmt.Printf("A player has requested a move %d\n")
	return &gm.Command{
		Command: "say hi",
	}, nil
}

func (s *GameService) RegisterPlayer(ctx context.Context, req *gm.Player) (*gm.Player, error) {
	fmt.Printf("New player with the name %s has entered the game\n", req.GetName())
	s.player = append(s.player, req.GetName())
	return req, nil
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
