package chess

func (gs *GameState) restoreGameState(ui UndoInfo) {
	gs.CastlingRights = ui.CastlingRights
	gs.HalfMoveClock = ui.HalfMoveClock
	gs.FullMoveCounter = ui.FullMoveCounter
	gs.EnPassantSquare = ui.EnPassantSquare
}

func (gs *GameState) restoreCapturedPiece(from, to Square, capturedPiece Piece) {
	board := gs.Board
	board[from] = board[to]
	board[to] = capturedPiece
}

func (gs *GameState) UnmakeNormalMove(m Move, ui UndoInfo) {
	board := gs.Board
	from := m.From
	to := m.To
	color := board[to].Color
	gs.restoreCapturedPiece(from, to, ui.CapturedPiece)
	if board[from].PieceType == King {
		gs.updateKingSquare(color, from)
	}
	gs.restoreGameState(ui)
}

func (gs *GameState) UnmakeDoublePush(m Move, ui UndoInfo) {
	from := m.From
	to := m.To
	gs.restoreCapturedPiece(from, to, ui.CapturedPiece)
	gs.restoreGameState(ui)
}

func (gs *GameState) UnmakeCastle(m Move, ui UndoInfo) {
	board := gs.Board
	from := m.From
	to := m.To
	board[ui.CastledRookFrom], board[ui.CastledRookTo] = board[ui.CastledRookTo], board[ui.CastledRookFrom]
	board[from], board[to] = board[to], board[from]
	color := board[from].Color
	gs.updateKingSquare(color, from)
	gs.restoreGameState(ui)
}

func (gs *GameState) UnmakeEnPassantSquare(m Move, ui UndoInfo) {
	board := gs.Board
	from := m.From
	to := m.To
	board[from], board[to] = board[to], board[from]
	board[ui.CapturedPawnSquare] = ui.CapturedPiece
	gs.restoreGameState(ui)
}

func (gs *GameState) UnmakePromotion(m Move, ui UndoInfo) {
	board := gs.Board
	from := m.From
	to := m.To
	color := board[to].Color
	board[from] = Piece{
		PieceType: Pawn,
		Color:     color,
	}
	board[to] = ui.CapturedPiece
	gs.restoreGameState(ui)
}
