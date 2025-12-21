package chess

import "fmt"

type Board [64]Piece

func NewBoard() Board {
	var board Board
	for i := range 64 {
		board[i] = NewPiece(PieceNone, ColorWhite)
	}
	return board
}

func (b *Board) Get(sq Square) Piece {
	return b[sq]
}
func (b *Board) Set(sq Square, p Piece) {
	if !sq.isValid() {
		return
	}
	b[sq] = p
}
func (b *Board) isEmpty(sq Square) bool {
	return b[sq].IsEmpty()
}
func (b *Board) isEnemy(sq Square, c Color) bool {
	if !b.isEmpty(sq) && b[sq].Color == c.Opponent() {
		return true
	}
	return false
}
func (b *Board) setPawnsOnRank(c Color) {
	var rank int
	if c == ColorWhite {
		rank = 1
	} else {
		rank = 6
	}

	for file := range 8 {
		sq := NewSquare(file, rank)
		b[int(sq)] = NewPiece(Pawn, c)
	}

}

func (b *Board) setPiecesOnRank(c Color) {
	pieces := [8]PieceType{Rook, Knight, Bishop, King, Queen, Bishop, Knight, Rook}
	var rank int
	if c == ColorWhite {
		rank = 0
	} else {
		rank = 7
	}
	for file := range 8 {
		sq := NewSquare(file, rank)
		b[int(sq)] = NewPiece(pieces[file], c)
	}
}

func (b *Board) DisplayBoard() {
	for rank := 7; rank >= 0; rank-- {
		fmt.Printf("%d | ", rank+1)
		for file := 0; file < 8; file++ {
			sq := Square(rank*8 + file)
			fmt.Print(b[sq].String(), " ")
		}
		fmt.Println()
	}

	fmt.Println("   ----------------")
	fmt.Print("    ")
	for file := 0; file < 8; file++ {
		fmt.Printf("%c ", 'a'+file)
	}
	fmt.Println()
}
