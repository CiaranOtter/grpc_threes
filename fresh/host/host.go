package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"three_engine/services/three"

	"google.golang.org/grpc"
)

type ThreeService struct {
	three.UnimplementedThreesServer
	players       []*three.Player
	turn          int32
	awaiting_move bool
	last_move     *three.Move
}

func (s *ThreeService) ConnectPlayer(ctx context.Context, req *three.Player) (*three.Player, error) {
	fmt.Printf("Player has requested to connect\n")

	req.Colour = int32(len(s.players))

	fmt.Printf("New player connected:\n\tName: %s\n\tColour: %d\n", req.GetName(), req.GetColour())

	s.players = append(s.players, req)
	s.turn = 0

	for !(len(s.players) == 2) {

	}
	s.awaiting_move = true
	return req, nil
}

func (s *ThreeService) GetCommand(ctx context.Context, in *three.Request) (*three.Response, error) {
	if in.PlayerId == s.turn {
		return &three.Response{
			Type: three.CMD_Type_MY_MOVE,
		}, nil
	} else {
		// wait for opponent to finish making their move
		for s.awaiting_move {

		}
		s.awaiting_move = true
		return &three.Response{
			Type:         three.CMD_Type_OPP_MOVE,
			PlayerColour: (in.GetPlayerId() + 1) % 2,
			Move:         s.last_move,
		}, nil
	}
}

func (s *ThreeService) Makemove(ctx context.Context, in *three.Request) (*three.Response, error) {
	s.turn = (s.turn + 1) % 2
	s.awaiting_move = false

	s.last_move = in.GetMove()

	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:5000")

	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption

	service := &ThreeService{}
	Server := grpc.NewServer(opts...)
	three.RegisterThreesServer(Server, service)

	Server.Serve(lis)
}
