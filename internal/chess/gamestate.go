package chess

type CastlingRights struct {
	WhiteKingSide  bool
	BlackKingSide  bool
	WhiteQueenSide bool
	BlackQueenSide bool
}

type GameState struct {
	Board           Board
	SideToMove      Color
	CastlingRights  CastlingRights
	EnPassantSquare Square
	HalfMoveClock   int
	FullMoveCounter int
	whiteKingCached Square
	blackKingCached Square
}

func (gs *GameState) GetKingSquare(c Color) Square {
	if c == ColorWhite {
		return gs.whiteKingCached
	}
	return gs.blackKingCached
}

func (gs *GameState) updateKingSquare(c Color, sq Square) {
	if c == ColorBlack {
		gs.blackKingCached = sq
	} else {
		gs.whiteKingCached = sq
	}
}

func (gs *GameState) Copy() GameState {
	gscopy := *gs
	return gscopy
}

func NewInitialGameState() GameState {
	board := NewBoard()
	board.setPawnsOnRank(ColorWhite)
	board.setPawnsOnRank(ColorBlack)
	board.setPiecesOnRank(ColorWhite)
	board.setPiecesOnRank(ColorBlack)

	return GameState{
		Board:      board,
		SideToMove: ColorWhite,
		CastlingRights: CastlingRights{
			WhiteKingSide:  true,
			BlackKingSide:  true,
			WhiteQueenSide: true,
			BlackQueenSide: true,
		},
		EnPassantSquare: Square(-1),
		HalfMoveClock:   0,
		FullMoveCounter: 1,
		whiteKingCached: NewSquare(4, 0),
		blackKingCached: NewSquare(4, 7),
	}
}
