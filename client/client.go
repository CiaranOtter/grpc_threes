package client

import (
	"context"
	"log"
	"time"

	"github.com/CiaranOtter/grpc_threes.git/gameservice"
	gm "github.com/CiaranOtter/grpc_threes.git/gameservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client_cnn gameservice.GameServiceClient

func MakeMove(move gm.Move) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client_cnn.MakeMove(ctx, &move)
	// defer cancel()
}

func GetCommand() *gameservice.Command {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	command, err := client_cnn.GetCommand(ctx, &gameservice.Empty{})

	if err != nil {
		log.Fatal(err)
	}

	// defer cancel()

	return command
}

func StartClient(address string, port int) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient("localhost:5000", opts...)

	if err != nil {
		log.Fatal(err)
	}

	// defer conn.Close()

	client_cnn = gm.NewGameServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client_cnn.RegisterPlayer(ctx, &gm.Player{
		Name: "Ciaran",
	})
	// defer cancel()
}
