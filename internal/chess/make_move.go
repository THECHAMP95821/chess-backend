package chess

func (gs *GameState) updateCastlingRights(from Square, to Square, capturedPiece Piece) {
	piece := gs.Board[from]
	color := piece.Color
	pieceType := piece.PieceType
	if pieceType == King {
		if color == ColorWhite && from == NewSquare(4, 0) {
			gs.CastlingRights.WhiteKingSide = false
			gs.CastlingRights.WhiteQueenSide = false
		} else if color == ColorBlack && from == NewSquare(4, 7) {
			gs.CastlingRights.BlackKingSide = false
			gs.CastlingRights.BlackQueenSide = false
		}
	}
	if pieceType == Rook {
		if color == ColorWhite && from == NewSquare(0, 0) {
			gs.CastlingRights.WhiteQueenSide = false
		} else if color == ColorWhite && from == NewSquare(7, 0) {
			gs.CastlingRights.WhiteKingSide = false
		} else if color == ColorBlack && from == NewSquare(0, 7) {
			gs.CastlingRights.BlackQueenSide = false
		} else if color == ColorBlack && from == NewSquare(7, 7) {
			gs.CastlingRights.BlackKingSide = false
		}
	}
	if capturedPiece.PieceType == Rook {
		if capturedPiece.Color == ColorWhite && to == NewSquare(0, 0) {
			gs.CastlingRights.WhiteQueenSide = false
		} else if capturedPiece.Color == ColorWhite && to == NewSquare(7, 0) {
			gs.CastlingRights.WhiteKingSide = false
		} else if capturedPiece.Color == ColorBlack && to == NewSquare(0, 7) {
			gs.CastlingRights.BlackQueenSide = false
		} else if capturedPiece.Color == ColorBlack && to == NewSquare(7, 7) {
			gs.CastlingRights.BlackKingSide = false
		}
	}
}

func (gs *GameState) updateEnPassantSquare(s Square, dir int) {
	enPassantSquare := s.applyOffset(dir * 8)
	gs.EnPassantSquare = enPassantSquare
}

func (gs *GameState) MakeNormalMove(m Move) UndoInfo {
	board := gs.Board
	from := m.From
	to := m.To
	piece := board[from]
	color := piece.Color
	pieceType := piece.PieceType
	capturedPiece := EmptyPiece()
	if m.IsCapture() {
		capturedPiece = board[to]
	}
	prevCastlingRights := gs.CastlingRights
	gs.updateCastlingRights(from, to, capturedPiece)

	board[to] = board[from]
	board[from] = EmptyPiece()

	if pieceType == King {
		gs.updateKingSquare(color, to)
	}

	return UndoInfo{
		CapturedPiece:      capturedPiece,
		CastlingRights:     prevCastlingRights,
		EnPassantSquare:    gs.EnPassantSquare,
		HalfMoveClock:      gs.HalfMoveClock,
		FullMoveCounter:    gs.FullMoveCounter,
		CastledRookFrom:    Square(-1),
		CastledRookTo:      Square(-1),
		CapturedPawnSquare: Square(-1),
	}
}

func (gs *GameState) MakeDoublePush(m Move) UndoInfo {
	from := m.From
	dir := -1
	if gs.Board[m.From].Color == ColorWhite {
		dir = 1
	}
	prevEnPassantSquare := gs.EnPassantSquare
	gs.updateEnPassantSquare(from, dir)
	return UndoInfo{
		CapturedPiece: Piece{
			PieceType: PieceNone,
			Color:     ColorWhite,
		},
		CastlingRights:     gs.CastlingRights,
		EnPassantSquare:    prevEnPassantSquare,
		HalfMoveClock:      gs.HalfMoveClock,
		FullMoveCounter:    gs.FullMoveCounter,
		CastledRookFrom:    Square(-1),
		CastledRookTo:      Square(-1),
		CapturedPawnSquare: Square(-1),
	}
}

func (gs *GameState) MakeCastle(m Move) UndoInfo {
	board := gs.Board
	from := m.From
	to := m.To
	color := board[from].Color
	prevCastlingRights := gs.CastlingRights
	gs.updateCastlingRights(from, to, EmptyPiece())
	board[to] = board[from]
	board[from] = EmptyPiece()
	gs.updateKingSquare(color, to)
	rookRank := 0
	if color == ColorBlack {
		rookRank = 7
	}
	initialRookSquare := Square(-1)
	finalRookSquare := Square(-1)
	if to.File() == FILEG {
		initialRookSquare = NewSquare(FILEH, rookRank)
		finalRookSquare = NewSquare(FILEF, rookRank)

	} else if to.File() == FILEC {
		initialRookSquare = NewSquare(FILEA, rookRank)
		finalRookSquare = NewSquare(FILED, rookRank)
	}
	board[finalRookSquare] = board[initialRookSquare]
	board[initialRookSquare] = EmptyPiece()
	return UndoInfo{
		CapturedPiece:      EmptyPiece(),
		CastlingRights:     prevCastlingRights,
		EnPassantSquare:    gs.EnPassantSquare,
		HalfMoveClock:      gs.HalfMoveClock,
		FullMoveCounter:    gs.FullMoveCounter,
		CastledRookFrom:    initialRookSquare,
		CastledRookTo:      finalRookSquare,
		CapturedPawnSquare: Square(-1),
	}
}

func (gs *GameState) MakeEnPassant(m Move) UndoInfo {
	board := gs.Board
	from := m.From
	to := m.To
	dir := -8 // White captures
	if gs.Board[from].Color == ColorBlack {
		dir = 8 // Black captures
	}
	capturedpawnSquare := Square(int(to) + dir)
	capturedPiece := board[capturedpawnSquare]
	board[capturedpawnSquare] = EmptyPiece()
	board[to] = board[from]
	board[from] = EmptyPiece()
	return UndoInfo{
		CapturedPiece:      capturedPiece,
		CastlingRights:     gs.CastlingRights,
		EnPassantSquare:    gs.EnPassantSquare,
		HalfMoveClock:      gs.HalfMoveClock,
		FullMoveCounter:    gs.FullMoveCounter,
		CastledRookFrom:    Square(-1),
		CastledRookTo:      Square(-1),
		CapturedPawnSquare: capturedpawnSquare,
	}
}

func (gs *GameState) MakePromotion(m Move) UndoInfo {
	board := gs.Board
	from := m.From
	to := m.To
	color := board[from].Color
	promoPieceType := m.Promotion
	capturedPiece := EmptyPiece()
	if m.IsCapture() {
		capturedPiece = board[to]
	}
	board[from] = EmptyPiece()
	board[to] = Piece{
		PieceType: promoPieceType,
		Color:     color,
	}
	return UndoInfo{
		CapturedPiece:      capturedPiece,
		CastlingRights:     gs.CastlingRights,
		EnPassantSquare:    gs.EnPassantSquare,
		HalfMoveClock:      gs.HalfMoveClock,
		FullMoveCounter:    gs.FullMoveCounter,
		CastledRookFrom:    Square(-1),
		CastledRookTo:      Square(-1),
		CapturedPawnSquare: Square(-1),
	}
}
