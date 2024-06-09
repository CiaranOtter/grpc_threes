package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"three_engine/services/three"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ClientConn three.ThreesClient
var running bool
var colour int32
var open_spaces []*Node
var num_pieces int

type Node struct {
	piece int32
	row   int32
	col   int32
}

var char_reps = map[int]string{
	-4: " ",
	-3: "|",
	-2: "-",
	-1: "O",
	0:  "W",
	1:  "B",
}

var board [7][7]*Node

func NewBoard() ([7][7]*Node, int) {
	return [7][7]*Node{
		{open_space(0, 0), hor_line(), hor_line(), open_space(0, 3), hor_line(), hor_line(), open_space(0, 6)},
		{vert_line(), open_space(1, 1), hor_line(), open_space(1, 3), hor_line(), open_space(1, 5), vert_line()},
		{vert_line(), vert_line(), open_space(2, 2), open_space(2, 3), open_space(2, 4), vert_line(), vert_line()},
		{open_space(3, 0), open_space(3, 1), open_space(3, 2), center(), open_space(3, 4), open_space(3, 5), open_space(3, 6)},
		{vert_line(), vert_line(), open_space(4, 2), open_space(4, 3), open_space(4, 4), vert_line(), vert_line()},
		{vert_line(), open_space(5, 1), hor_line(), open_space(5, 3), hor_line(), open_space(5, 5), vert_line()},
		{open_space(6, 0), hor_line(), hor_line(), open_space(6, 3), hor_line(), hor_line(), open_space(6, 6)},
	}, 8
}

func vert_line() *Node {
	return &Node{
		piece: -3,
	}
}

func hor_line() *Node {
	return &Node{
		piece: -2,
	}
}

func open_space(row int, col int) *Node {
	return &Node{
		piece: -1,
		row:   int32(row),
		col:   int32(col),
	}
}

func center() *Node {
	return &Node{
		piece: -4,
	}
}

func printBoard() {
	for _, row := range board {
		for _, col := range row {
			fmt.Printf(" %s", char_reps[int(col.piece)])
		}
		fmt.Printf("\n")
	}
}

func find_move() {
	open_spaces = make([]*Node, 0)
	for _, row := range board {
		for _, col := range row {
			if col.piece == -1 {
				open_spaces = append(open_spaces, col)
			}
		}
	}
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient("localhost:5000", opts...)

	if err != nil {
		log.Fatal(err)
	}

	board, num_pieces = NewBoard()
	printBoard()

	ClientConn = three.NewThreesClient(conn)

	resp, err := ClientConn.ConnectPlayer(context.Background(), &three.Player{
		Name: "Ciaran",
	})

	if err != nil {
		conn.Close()
		log.Fatal(err)
	}

	fmt.Printf("Player has been registered: \n\tName: %s\n\tColour: %d\n", resp.GetName(), resp.GetColour())

	colour = resp.GetColour()
	running = true

	play_loop()
}

func play_loop() {
	for running {

		time.Sleep(1 * time.Second)
		command, err := ClientConn.GetCommand(context.Background(), &three.Request{
			PlayerId: colour,
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Received command of type %s from server\n", command.GetType().String())
		handleCommand(command)
	}
}

func handleCommand(command *three.Response) {
	switch command.GetType() {
	case (three.CMD_Type_MY_MOVE):
		make_move()
		break
	case (three.CMD_Type_OPP_MOVE):
		add_opp_move(command)
		break
	default:
		fmt.Printf("Unknown Move\n")
		running = false
		break
	}
}

func apply_move(row int, col int, colour int32) {
	board[row][col].piece = colour
}

func add_opp_move(command *three.Response) {
	fmt.Printf("Making my opponents move\n")
	apply_move(int(command.GetMove().GetRow()), int(command.GetMove().GetCol()), command.GetPlayerColour())
}

func make_move() {
	fmt.Printf("Making my move\n")

	find_move()

	chosenMove := rand.Intn(len(open_spaces))

	_, err := ClientConn.Makemove(context.Background(), &three.Request{
		PlayerId: colour,
		Move: &three.Move{
			Type: three.Type_ADD,
			Row:  open_spaces[chosenMove].row,
			Col:  open_spaces[chosenMove].col,
		},
	})

	apply_move(int(open_spaces[chosenMove].row), int(open_spaces[chosenMove].col), colour)

	num_pieces--

	printBoard()

	if err != nil {
		log.Fatal(err)
	}
}
