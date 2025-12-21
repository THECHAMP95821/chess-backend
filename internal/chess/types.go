package chess

import (
	"fmt"
	"strings"
)

type Color uint8

const (
	ColorWhite Color = iota
	ColorBlack
)

func (c Color) Opponent() Color {
	return Color(c ^ 1)
}

func (c Color) String() string {
	switch c {
	case ColorWhite:
		return "White"
	case ColorBlack:
		return "Black"
	default:
		return "INVALID COLOR"
	}
}

type PieceType uint8

const (
	PieceNone PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

func (p PieceType) String() string {
	switch p {
	case PieceNone:
		return "."
	case Pawn:
		return "P"
	case Knight:
		return "N"
	case Bishop:
		return "B"
	case Rook:
		return "R"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return "INVALID PIECE"
	}
}

func ParsePieceType(s string) PieceType {
	if s == "." {
		return PieceNone
	}

	switch strings.ToUpper(s) {
	case "P":
		return Pawn
	case "N":
		return Knight
	case "B":
		return Bishop
	case "R":
		return Rook
	case "Q":
		return Queen
	case "K":
		return King
	default:
		return PieceNone
	}
}

type Piece struct {
	PieceType PieceType
	Color     Color
}

func NewPiece(pt PieceType, c Color) Piece {
	return Piece{
		PieceType: pt,
		Color:     c,
	}
}

func (p Piece) Type() PieceType {
	return p.PieceType
}

func (p Piece) getColor() Color {
	return p.Color
}
func (p Piece) IsEmpty() bool {
	return p.PieceType == PieceNone
}
func (p Piece) IsOpponent(piece Piece) bool {
	return !p.IsEmpty() && p.Color == piece.Color.Opponent()
}
func (p Piece) IsAlly(piece Piece) bool {
	return !p.IsEmpty() && p.Color == piece.Color
}
func (p Piece) String() string {
	if p.IsEmpty() {
		return p.PieceType.String()
	}
	if p.Color == ColorBlack {
		return strings.ToLower(p.PieceType.String())
	}
	return p.PieceType.String()
}

type Square int8

func NewSquare(file, rank int) Square {
	return Square(file + rank*8)
}

func (s Square) File() int {
	return int(s) % 8
}
func (s Square) Rank() int {
	return int(s) / 8
}
func (s Square) isValid() bool {
	return int(s) >= 0 && int(s) < 64
}
func (s Square) String() string {
	return fmt.Sprintf("%c%d", 'a'+s.File(), s.Rank()+1)
}
func (s Square) applyOffset(offset int) Square {
	return Square(int(s) + offset)
}
func ParseSquare(s string) Square {
	file := int(s[0] - 'a')
	rank := int(s[1] - '1')
	return NewSquare(file, rank)
}

var KnightOffsets = [8]int{-17, -15, -10, -6, 6, 10, 15, 17}
var KingOffsets = [8]int{-9, -8, -7, -1, 1, 7, 8, 9}
var BishopOffsets = [4]int{-9, -7, 7, 9}
var RookOffsets = [4]int{-8, -1, 1, 8}
var QueenOffsets = [8]int{-9, -8, -7, -1, 1, 7, 8, 9}
