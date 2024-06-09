package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	gm "github.com/CiaranOtter/grpc_threes.git/gameservice"
	"google.golang.org/grpc"
)

type GameService struct {
	gm.UnimplementedGameServiceServer
	// players
	player         []player
	current_player int
}

type player struct {
	index       int
	name        string
	played_move bool
	colour      int
	last_move   *gm.Move
}

func (s *GameService) MakeMove(ctx context.Context, req *gm.Move) (*gm.Response, error) {
	fmt.Printf("Player with the colour %d made move %d, %d\n", req.Colour, req.Row, req.Col)

	s.player[req.GetRequest().GetPlayerIndex()].played_move = true
	s.player[req.GetRequest().GetPlayerIndex()].last_move = req

	s.current_player = (s.current_player + 1) % 2

	return &gm.Response{
		Success: true,
	}, nil
}

func (s *GameService) GetCommand(ctx context.Context, req *gm.Request) (*gm.Command, error) {
	fmt.Printf("A player has requested a move %d\n")

	if req.GetPlayerIndex() == int32(s.current_player) {

		return &gm.Command{
			Command: gm.CMD_MAKE_MOVE,
		}, nil
	} else {
		for !s.player[(req.GetPlayerIndex()+1)%2].played_move {
			fmt.Printf("Awaiting opponent move")
			time.Sleep(500 * time.Millisecond)
		}
		return &gm.Command{
			Command: gm.CMD_PLAY_MOVE,
			Move:    s.player[(s.current_player+1)%2].last_move,
		}, nil
	}

}

func (s *GameService) RegisterPlayer(ctx context.Context, req *gm.Player) (*gm.Player, error) {
	fmt.Printf("New player with the name %s has entered the game\n", req.GetName())

	s.player = append(s.player, player{
		index:       len(s.player) - 1,
		colour:      len(s.player) - 1,
		name:        req.Name,
		played_move: false,
	})

	req.Index = int32(len(s.player) - 1)
	req.Colour = int32(len(s.player) - 1)

	for _, player := range s.player {
		fmt.Printf("player: %s\n", player.name)
	}

	for len(s.player) != 2 {
		fmt.Printf("awaiting another player...")
		time.Sleep(500 * time.Millisecond)
	}

	s.current_player = 0

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
