package chess

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func GeneratePawnMoves(state *GameState, from Square, moves *[]Move) {
	board := state.Board
	color := board[from].Color
	dir := -1
	if color == ColorWhite {
		dir = 1
	}
	flag := MoveFlagNone

	if (color == ColorWhite && from.Rank() == 6) || (color == ColorBlack && from.Rank() == 1) {
		flag = MoveFlagPromotion
	}

	singlePushSquare := Square(int(from) + dir*8)
	if singlePushSquare.isValid() && board[singlePushSquare].IsEmpty() {
		*moves = append(*moves, Move{
			From:  from,
			To:    singlePushSquare,
			Flags: flag,
		})
	}
	if (color == ColorWhite && from.Rank() == 1) || (color == ColorBlack && from.Rank() == 6) {
		doublePushSquare := Square(int(from) + 2*dir*8)
		if doublePushSquare.isValid() && board[doublePushSquare].IsEmpty() && board[singlePushSquare].IsEmpty() {
			*moves = append(*moves, Move{
				From:  from,
				To:    doublePushSquare,
				Flags: MoveFlagDoublePush,
			})
		}
	}
	captureOffSet := []int{}
	if color == ColorWhite {
		captureOffSet = []int{7, 9}
	} else {
		captureOffSet = []int{-7, -9}
	}
	for _, offset := range captureOffSet {
		captureSquare := from.applyOffset(offset)
		if captureSquare.isValid() && board[captureSquare].IsOpponent(board[from]) &&
			abs(from.File()-captureSquare.File()) == 1 {
			if captureSquare == state.EnPassantSquare {
				flag = MoveFlagEnPassant
			}
			*moves = append(*moves, Move{
				From:  from,
				To:    captureSquare,
				Flags: flag | MoveFlagCapture,
			})
		}
	}
}

func GenerateKnightMoves(state *GameState, from Square, moves *[]Move) {
	board := state.Board
	for _, offset := range KnightOffsets {
		to := from.applyOffset(offset)
		if !to.isValid() {
			continue
		}
		if !(abs(to.File()-from.File()) <= 2 && abs(to.Rank()-from.Rank()) <= 2) {
			continue
		}
		if board[to].IsAlly(board[from]) {
			continue
		}
		flag := MoveFlagNone
		if board[to].IsOpponent(board[from]) {
			flag = MoveFlagCapture
		}
		*moves = append(*moves, Move{
			From:  from,
			To:    to,
			Flags: flag,
		})
	}
}

func checkValidRankDiff(from, to Square) bool {
	return abs(from.File()-to.File()) <= 1 && abs(from.Rank()-to.Rank()) <= 1
}

func GenerateSlidingMoves(state *GameState, from Square, moves *[]Move, offset int) {
	board := state.Board
	i := 1
	stopMoves := false
	for !stopMoves {
		moveTo := from.applyOffset(offset)
		if !moveTo.isValid() {
			break
		}
		if !checkValidRankDiff(from, moveTo) {
			break
		}
		if board[moveTo].IsEmpty() {
			*moves = append(*moves, Move{
				From: from,
				To:   moveTo,
			})
		} else {
			if board[moveTo].IsOpponent(board[from]) {
				*moves = append(*moves, Move{
					From:  from,
					To:    moveTo,
					Flags: MoveFlagCapture,
				})
			}
			stopMoves = true
		}
		from = moveTo
		i++
	}
}

func GenerateBishopMoves(state *GameState, from Square, moves *[]Move) {
	for _, offset := range BishopOffsets {
		GenerateSlidingMoves(state, from, moves, offset)
	}
}

func GenerateQueenMoves(state *GameState, from Square, moves *[]Move) {
	for _, offset := range QueenOffsets {
		GenerateSlidingMoves(state, from, moves, offset)
	}
}

func GenerateRookMoves(state *GameState, from Square, moves *[]Move) {
	for _, offset := range RookOffsets {
		GenerateSlidingMoves(state, from, moves, offset)
	}
}

func isSquareSafeForKing(state *GameState, sq Square, color Color) bool {
	return !IsSquareAttacked(state, sq, color.Opponent())
}

func GenerateKingMoves(state *GameState, from Square, moves *[]Move) {
	board := state.Board
	for _, offset := range KingOffsets {
		to := from.applyOffset(offset)
		if !to.isValid() {
			continue
		}
		if !checkValidRankDiff(from, to) {
			continue
		}
		if board[to].IsAlly(board[from]) {
			continue
		}
		flag := MoveFlagNone
		if board[to].IsOpponent(board[from]) {
			flag = MoveFlagCapture
		}
		*moves = append(*moves, Move{
			From:  from,
			To:    to,
			Flags: flag,
		})
	}

	color := board[from].Color
	rank := 0
	if color == ColorBlack {
		rank = 7
	}

	if from != NewSquare(4, rank) {
		return
	}

	if IsSquareAttacked(state, from, color.Opponent()) {
		return
	}

	if (color == ColorWhite && state.CastlingRights.WhiteKingSide) ||
		(color == ColorBlack && state.CastlingRights.BlackKingSide) {

		f := NewSquare(5, rank)
		g := NewSquare(6, rank)

		if board[f].IsEmpty() &&
			board[g].IsEmpty() &&
			isSquareSafeForKing(state, f, color) &&
			isSquareSafeForKing(state, g, color) {

			*moves = append(*moves, Move{
				From:  from,
				To:    g,
				Flags: MoveFlagCastle,
			})
		}
	}

	if (color == ColorWhite && state.CastlingRights.WhiteQueenSide) ||
		(color == ColorBlack && state.CastlingRights.BlackQueenSide) {

		d := NewSquare(3, rank)
		c := NewSquare(2, rank)
		b := NewSquare(1, rank)

		if board[d].IsEmpty() &&
			board[c].IsEmpty() &&
			board[b].IsEmpty() &&
			isSquareSafeForKing(state, d, color) &&
			isSquareSafeForKing(state, c, color) {

			*moves = append(*moves, Move{
				From:  from,
				To:    c,
				Flags: MoveFlagCastle,
			})
		}
	}

}

func GeneratePseudoLegalMoves(state *GameState) []Move {
	board := state.Board
	moves := []Move{}
	for i := range 63 {
		sq := Square(i)
		piece := board[sq]
		if piece.Color == state.SideToMove && !piece.IsEmpty() {
			switch piece.PieceType {
			case Pawn:
				GeneratePawnMoves(state, sq, &moves)
			case Knight:
				GenerateKnightMoves(state, sq, &moves)
			case Bishop:
				GenerateBishopMoves(state, sq, &moves)
			case Rook:
				GenerateRookMoves(state, sq, &moves)
			case Queen:
				GenerateQueenMoves(state, sq, &moves)
			case King:
				GenerateKingMoves(state, sq, &moves)
			default:
				continue
			}

		}
	}
	return moves
}

func isSquareAttackedByPawn(state *GameState, square Square, byColor Color) bool {
	board := state.Board
	var captureOffsets []int
	if byColor == ColorWhite {
		captureOffsets = []int{-7, -9}
	} else {
		captureOffsets = []int{7, 9}
	}
	for _, offset := range captureOffsets {
		from := square.applyOffset(offset)
		if board[from].PieceType == Pawn && board[from].Color == byColor && checkValidRankDiff(from, square) {
			return true
		}
	}
	return false
}

func isSquareAttackedByKnight(state *GameState, square Square, byColor Color) bool {
	board := state.Board
	for _, offset := range KnightOffsets {
		from := square.applyOffset(offset)
		if board[from].PieceType == Knight && board[from].Color == byColor && abs(square.File()-from.File()) <= 2 && abs(square.Rank()-from.Rank()) <= 2 {
			return true
		}
	}
	return false
}

func isSquareAttackedByKing(state *GameState, square Square, byColor Color) bool {
	board := state.Board
	for _, offset := range KingOffsets {
		from := square.applyOffset(offset)
		if board[from].PieceType == King && board[from].Color == byColor && checkValidRankDiff(from, square) {
			return true
		}
	}
	return false
}

func isSquareAttackedSliding(state *GameState, from Square, offset int, byColor Color, pieceType PieceType) bool {
	board := state.Board
	i := 1
	stopMoves := false
	for !stopMoves {
		moveTo := from.applyOffset(offset)
		if !moveTo.isValid() {
			break
		}
		if !checkValidRankDiff(from, moveTo) {
			break
		}
		if !board[moveTo].IsEmpty() {
			if board[moveTo].Color == byColor {
				if board[moveTo].PieceType == pieceType {
					return true
				} else {
					return false
				}
			}
			stopMoves = true
		}
		from = moveTo
		i++
	}
	return false
}

func isSquareAttackedByBishop(state *GameState, square Square, byColor Color) bool {
	for _, offset := range BishopOffsets {
		if isSquareAttackedSliding(state, square, offset, byColor, Bishop) {
			return true
		}
	}
	return false
}

func isSquareAttackedByQueen(state *GameState, square Square, byColor Color) bool {
	for _, offset := range QueenOffsets {
		if isSquareAttackedSliding(state, square, offset, byColor, Queen) {
			return true
		}
	}
	return false
}

func isSquareAttackedByRook(state *GameState, square Square, byColor Color) bool {
	for _, offset := range RookOffsets {
		if isSquareAttackedSliding(state, square, offset, byColor, Rook) {
			return true
		}
	}
	return false
}

func IsSquareAttacked(state *GameState, square Square, byColor Color) bool {
	if !square.isValid() {
		return false
	}
	return isSquareAttackedByPawn(state, square, byColor) ||
		isSquareAttackedByKnight(state, square, byColor) ||
		isSquareAttackedByBishop(state, square, byColor) ||
		isSquareAttackedByRook(state, square, byColor) ||
		isSquareAttackedByQueen(state, square, byColor) ||
		isSquareAttackedByKing(state, square, byColor)
}

func (gs *GameState) IsKingInCheck() bool {
	var kingPos Square
	color := gs.SideToMove
	if color == ColorWhite {
		kingPos = gs.whiteKingCached
	} else {
		kingPos = gs.blackKingCached
	}
	return IsSquareAttacked(gs, kingPos, color.Opponent())
}

func GenerateLegalMoves(state *GameState) []Move {
	pseudoLegalMoves := GeneratePseudoLegalMoves(state)
	legalMoves := make([]Move, 0, len(pseudoLegalMoves))

	for _, move := range pseudoLegalMoves {
		undo := state.MakeMove(move)

		if !state.IsKingInCheck() {
			legalMoves = append(legalMoves, move)
		}

		state.UnmakeMove(move, undo)
	}
	return legalMoves
}
