package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/CiaranOtter/grpc_threes.git/client"
	gm "github.com/CiaranOtter/grpc_threes.git/gameservice"
	"github.com/CiaranOtter/grpc_threes.git/util"
	// "service/gameservice"
)

var board [7][7]*node
var moves = make([]*node, 0)
var open_space = make([]gm.Move, 0)

/** reset the game state
 *
 * @return:
 	- The game board as an array
	- The number of pieces
	- The phase
 *
*/

var char_reps = map[int]string{
	-3: " ",
	-2: "|",
	-1: "-",
	0:  "O",
	1:  "W",
	2:  "B",
}

/*
  - a node value for a point on the board
    space_value - the value of the space
    neighbors - the adjacent spaces to the node

*/

type node struct {
	space_value int
	piece       *util.Piece
	position    gm.Move
	neighbors   []*node
}

func (n *node) add_neighbor(neigh *node) *node {
	n.neighbors = append(n.neighbors, neigh)
	neigh.neighbors = append(neigh.neighbors, n)
	return neigh
}

func newBoard() ([7][7]*node, int, int) {
	return [7][7]*node{
		{space(0, 0), hor_line(), hor_line(), space(0, 3), hor_line(), hor_line(), space(0, 6)},
		{vert_line(), space(1, 1), hor_line(), space(1, 3), hor_line(), space(1, 5), vert_line()},
		{vert_line(), vert_line(), space(2, 2), space(2, 3), space(2, 4), vert_line(), vert_line()},
		{space(3, 0), space(3, 1), space(3, 2), center(), space(3, 4), space(3, 5), space(3, 6)},
		{vert_line(), vert_line(), space(4, 2), space(4, 3), space(4, 4), vert_line(), vert_line()},
		{vert_line(), space(5, 1), hor_line(), space(5, 3), hor_line(), space(5, 5), vert_line()},
		{space(6, 0), hor_line(), hor_line(), space(6, 3), hor_line(), hor_line(), space(6, 6)},
	}, 8, 1
}

/*
* returns a center point on the board as a node
 */
func center() *node {
	return &node{
		space_value: -3,
		neighbors:   nil,
	}
}

/** returns a vertical line as a node */
func vert_line() *node {
	return &node{
		space_value: -2,
		neighbors:   nil,
	}
}

/*
* returns a horixontal line as a node
 */
func hor_line() *node {
	return &node{
		space_value: -1,
		neighbors:   nil,
	}
}

func space(row int, col int) *node {
	space := &node{
		space_value: 0,
		position: gm.Move{
			Row: int32(row),
			Col: int32(col),
		},
		neighbors: make([]*node, 0),
	}

	moves = append(moves, space)
	return space
}

var num_pieces int
var phase int

var pieces []*util.Piece

/** evalutate the board to generate score for player
 */
func evaluate() {
	// for _, row := range board {
	// 	for _, cell := range board[i] {
	// 		if board[i][j] == -1 {
	// 			continue
	// 		}
	// 		fmt.Printf("Evaluating col %d\n", j)
	// 	}
	// }
}

func print_board(b [7][7]*node) {
	for _, row := range b {
		for _, cell := range row {
			if cell.piece != nil {
				fmt.Printf("%s", char_reps[cell.piece.Colour])
			} else {
				fmt.Printf("%s", char_reps[cell.space_value])
			}
		}
		fmt.Printf("\n")
	}
}

func play() {
	switch phase {
	case 1:
		break
	case 2:
		break
	case 3:
		break
	}
}

func connect_board() {
	board[0][0].add_neighbor(board[0][3]).add_neighbor(board[0][6]).add_neighbor(board[3][6]).add_neighbor(board[6][6])
	board[0][0].add_neighbor(board[3][0]).add_neighbor(board[6][0]).add_neighbor(board[6][3]).add_neighbor(board[6][6])

	board[1][1].add_neighbor(board[1][3]).add_neighbor(board[1][5]).add_neighbor(board[3][5]).add_neighbor(board[5][5])
	board[1][1].add_neighbor(board[3][1]).add_neighbor(board[5][1]).add_neighbor(board[5][3]).add_neighbor(board[5][5])

	board[2][2].add_neighbor(board[2][3]).add_neighbor(board[2][4]).add_neighbor(board[3][4]).add_neighbor(board[4][4])
	board[2][2].add_neighbor(board[3][2]).add_neighbor(board[4][2]).add_neighbor(board[4][3]).add_neighbor(board[4][4])

	board[3][0].add_neighbor(board[3][1]).add_neighbor(board[3][2])
	board[3][4].add_neighbor(board[3][5]).add_neighbor(board[3][6])

	board[0][3].add_neighbor(board[1][3]).add_neighbor(board[2][3])
	board[4][3].add_neighbor(board[5][3]).add_neighbor(board[6][3])
}

func find_move() {
	open_space = make([]gm.Move, 0)
	for i := range board {
		for j := range board {
			if board[i][j].space_value == 0 && board[i][j].piece == nil {
				open_space = append(open_space, gm.Move{
					Row: int32(i),
					Col: int32(j),
				})
			}
		}
	}
}

func print_connections() {
	for _, node := range moves {
		fmt.Printf("Analysing position %d, %d: ", node.position.Row, node.position.Col)
		for _, n := range node.neighbors {
			fmt.Printf(" (%d, %d) ->", n.position.Row, n.position.Col)
		}
		fmt.Printf("\n")
	}
}

func add_piece(board *[7][7]*node, position *node, colour int) {
	num_pieces--
	fmt.Printf("Setting row %d, ol %d \n", position.position.Row, position.position.Col)

	p := util.Piece{
		Position: position.position,
		Safe:     false,
		Colour:   colour,
	}

	board[position.position.Row][position.position.Col].piece = &p

	pieces = append(pieces, &p)
}

func move_piece(board *[7][7]*node, piece *util.Piece) {
	fmt.Printf("Moving piece: (%d, %d)", piece.Position.Row, piece.Position.Col)
	for _, node := range board[piece.Position.Row][piece.Position.Col].neighbors {
		if node.piece == nil {
			fmt.Printf(" -> (%d, %d)", node.position.Row, node.position.Col)
			fmt.Printf("\n")
			board[piece.Position.Row][piece.Position.Col].piece = nil
			piece.Position = node.position
			node.piece = piece
			break
		}
	}
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Connect to port: ")
	port, _ := reader.ReadString('\n')
	port_num, err := strconv.Atoi(port)

	if err == nil {
		log.Fatal(err)
	}
	client.StartClient("localhost", int(port_num))

	board, num_pieces, phase = newBoard()
	connect_board()

	print_connections()
	print_board(board)

	for num_pieces > 0 {
		find_move()
		chosenMove := rand.Intn(len(open_space))

		// fmt.Printf("Chose move %d out of %d moves\n", chosenMove, len(moves))

		add_piece(&board, board[open_space[chosenMove].Row][open_space[chosenMove].Col], 1)
		print_board(board)
		fmt.Printf("\n")

	}

	for i := 0; i < 10; i++ {
		chosenPiece := rand.Intn(8)

		move_piece(&board, pieces[chosenPiece])
		print_board(board)
		fmt.Printf("\n")
	}
	// for i := 0; i < 10; i++ {
	// 	find_move(&board)

	// 	chosenMove := rand.Intn(len(moves))
	// 	chosenPiece := rand.Intn(len(pieces))

	// 	move_piece(&board, &pieces[chosenPiece], moves[chosenMove])

	// 	print_board(board)
	// 	fmt.Printf("\n")
	// }
	// evaluate()
}

func minimax() {

}
