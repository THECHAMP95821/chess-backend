package chess

import (
	"fmt"
	"testing"
)

func TestMoveGen(t *testing.T) {
	var board Board
	sq := ParseSquare("d5")
	fmt.Println(sq)
	board[28] = Piece{
		PieceType: King,
		Color:     ColorWhite,
	}
	board[35] = Piece{
		PieceType: Pawn,
		Color:     ColorBlack,
	}
	board[31] = Piece{
		PieceType: Rook,
		Color:     ColorBlack,
	}
	board[1] = Piece{
		PieceType: Queen,
		Color:     ColorBlack,
	}
	board[63] = Piece{
		PieceType: Bishop,
		Color:     ColorBlack,
	}
	board[29] = Piece{
		PieceType: Bishop,
		Color:     ColorBlack,
	}
	board[49] = Piece{
		PieceType: Queen,
		Color:     ColorBlack,
	}
	//board[55] = Piece{
	//	PieceType: Queen,
	//	Color:     ColorBlack,
	//}
	board[43] = Piece{
		PieceType: Knight,
		Color:     ColorBlack,
	}
	state := GameState{
		Board: board,
	}
	board.DisplayBoard()
	a := isSquareAttackedByPawn(&state, Square(28), ColorBlack)
	a = isSquareAttackedByKnight(&state, Square(28), ColorBlack)
	a = isSquareAttackedByBishop(&state, Square(28), ColorBlack)
	a = isSquareAttackedByQueen(&state, Square(28), ColorBlack)
	a = isSquareAttackedByRook(&state, Square(28), ColorBlack)
	fmt.Println(a)
}
