package chess

import (
	"errors"
	"strconv"
	"strings"
)

func pieceFromFEN(c rune) Piece {
	color := ColorWhite
	if c >= 'a' && c <= 'z' {
		color = ColorBlack
		c = rune(c - 32)
	}

	switch c {
	case 'P':
		return Piece{Pawn, color}
	case 'N':
		return Piece{Knight, color}
	case 'B':
		return Piece{Bishop, color}
	case 'R':
		return Piece{Rook, color}
	case 'Q':
		return Piece{Queen, color}
	case 'K':
		return Piece{King, color}
	default:
		return EmptyPiece()
	}
}

func pieceToFEN(p Piece) string {
	if p.IsEmpty() {
		return ""
	}

	var c rune
	switch p.PieceType {
	case Pawn:
		c = 'P'
	case Knight:
		c = 'N'
	case Bishop:
		c = 'B'
	case Rook:
		c = 'R'
	case Queen:
		c = 'Q'
	case King:
		c = 'K'
	}

	if p.Color == ColorBlack {
		c += 32
	}

	return string(c)
}

func ParseFEN(fen string) (*GameState, error) {
	parts := strings.Split(fen, " ")
	if len(parts) != 6 {
		return nil, errors.New("invalid FEN")
	}

	board := Board{}
	ranks := strings.Split(parts[0], "/")
	if len(ranks) != 8 {
		return nil, errors.New("invalid board in FEN")
	}

	sq := Square(56) // a8
	for _, rank := range ranks {
		file := 0
		for _, ch := range rank {
			if ch >= '1' && ch <= '8' {
				file += int(ch - '0')
			} else {
				board[sq+Square(file)] = pieceFromFEN(ch)
				file++
			}
		}
		sq -= 8
	}

	side := ColorWhite
	if parts[1] == "b" {
		side = ColorBlack
	}

	cr := CastlingRights{}
	if parts[2] != "-" {
		cr.WhiteKingSide = strings.Contains(parts[2], "K")
		cr.WhiteQueenSide = strings.Contains(parts[2], "Q")
		cr.BlackKingSide = strings.Contains(parts[2], "k")
		cr.BlackQueenSide = strings.Contains(parts[2], "q")
	}

	ep := Square(-1)
	if parts[3] != "-" {
		file := int(parts[3][0] - 'a')
		rank := int(parts[3][1] - '1')
		ep = NewSquare(file, rank)
	}

	halfMove, _ := strconv.Atoi(parts[4])
	fullMove, _ := strconv.Atoi(parts[5])

	state := &GameState{
		Board:           board,
		SideToMove:      side,
		CastlingRights:  cr,
		EnPassantSquare: ep,
		HalfMoveClock:   halfMove,
		FullMoveCounter: fullMove,
	}

	state.cacheKingSquares()
	return state, nil
}

func (gs *GameState) ToFEN() string {
	var sb strings.Builder

	for rank := 7; rank >= 0; rank-- {
		empty := 0
		for file := 0; file < 8; file++ {
			p := gs.Board[NewSquare(file, rank)]
			if p.IsEmpty() {
				empty++
			} else {
				if empty > 0 {
					sb.WriteString(strconv.Itoa(empty))
					empty = 0
				}
				sb.WriteString(pieceToFEN(p))
			}
		}
		if empty > 0 {
			sb.WriteString(strconv.Itoa(empty))
		}
		if rank != 0 {
			sb.WriteRune('/')
		}
	}

	if gs.SideToMove == ColorWhite {
		sb.WriteString(" w ")
	} else {
		sb.WriteString(" b ")
	}

	cr := ""
	if gs.CastlingRights.WhiteKingSide {
		cr += "K"
	}
	if gs.CastlingRights.WhiteQueenSide {
		cr += "Q"
	}
	if gs.CastlingRights.BlackKingSide {
		cr += "k"
	}
	if gs.CastlingRights.BlackQueenSide {
		cr += "q"
	}
	if cr == "" {
		cr = "-"
	}
	sb.WriteString(cr + " ")

	if gs.EnPassantSquare == -1 {
		sb.WriteString("- ")
	} else {
		file := rune('a' + gs.EnPassantSquare.File())
		rank := rune('1' + gs.EnPassantSquare.Rank())
		sb.WriteRune(file)
		sb.WriteRune(rank)
		sb.WriteString(" ")
	}

	sb.WriteString(strconv.Itoa(gs.HalfMoveClock))
	sb.WriteRune(' ')
	sb.WriteString(strconv.Itoa(gs.FullMoveCounter))

	return sb.String()
}
