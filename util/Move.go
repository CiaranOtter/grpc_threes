package util

import (
	gm "github.com/CiaranOtter/grpc_threes.git/gameservice"
)

/*
  - A piece on the board
    position - a move for the position of the piec
    safe - the value for whether a piece is protected or not
    colour - the colour of the piece
    poss_move - the possible moves for the piece
*/
type Piece struct {
	Position  gm.Move
	Safe      bool
	Colour    int
	Poss_move []int
}
