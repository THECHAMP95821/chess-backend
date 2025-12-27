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

func (gs *GameState) cacheKingSquares() {
	for sq := Square(0); sq < 64; sq++ {
		piece := gs.Board[sq]
		if piece.PieceType != King {
			continue
		}

		if piece.Color == ColorWhite {
			gs.whiteKingCached = sq
		} else {
			gs.blackKingCached = sq
		}
	}
}

func (gs *GameState) Copy() GameState {
	gscopy := *gs
	return gscopy
}
func (gs *GameState) switchSides() {
	gs.SideToMove = gs.SideToMove.Opponent()
}

func (gs *GameState) updateClocks(m Move) {
	from := m.From
	piece := gs.Board[from]
	if piece.Color == ColorBlack {
		gs.FullMoveCounter++
	}
	if piece.PieceType == Pawn || m.IsCapture() {
		gs.HalfMoveClock = 0
	} else {
		gs.HalfMoveClock++
	}
}

func (gs *GameState) MakeMove(m Move) UndoInfo {
	undoInfo := UndoInfo{}
	switch m.Flags {
	case MoveFlagNone, MoveFlagCapture:
		undoInfo = gs.MakeNormalMove(m)
	case MoveFlagDoublePush:
		undoInfo = gs.MakeDoublePush(m)
	case MoveFlagCastle:
		undoInfo = gs.MakeCastle(m)
	case MoveFlagEnPassant:
		undoInfo = gs.MakeEnPassant(m)
	case MoveFlagPromotion:
		undoInfo = gs.MakePromotion(m)
	}
	if m.Flags != MoveFlagDoublePush {
		gs.EnPassantSquare = Square(-1)
	}
	gs.switchSides()
	gs.updateClocks(m)
	return undoInfo
}

func (gs *GameState) UnmakeMove(m Move, ui UndoInfo) {
	gs.switchSides()
	switch m.Flags {
	case MoveFlagNone, MoveFlagCapture:
		gs.UnmakeNormalMove(m, ui)
	case MoveFlagDoublePush:
		gs.UnmakeDoublePush(m, ui)
	case MoveFlagCastle:
		gs.UnmakeCastle(m, ui)
	case MoveFlagEnPassant:
		gs.UnmakeEnPassantSquare(m, ui)
	case MoveFlagPromotion:
		gs.UnmakePromotion(m, ui)
	}
}
