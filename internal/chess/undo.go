package chess

type UndoInfo struct {
	CapturedPiece      Piece
	CastlingRights     CastlingRights
	EnPassantSquare    Square
	HalfMoveClock      int
	FullMoveCounter    int
	CastledRookFrom    Square
	CastledRookTo      Square
	CapturedPawnSquare Square
}
