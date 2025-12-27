package chess

type GameResult int

const (
	GameOngoing GameResult = iota
	GameWhiteWins
	GameBlackWins
	GameDraw
)

type DrawReason int

const (
	DrawNone DrawReason = iota
	DrawStalemate
	DrawInsufficientMaterial
	DrawFiftyMoveRule
	DrawThreefoldRepetition
)

type Outcome struct {
	Result     GameResult
	DrawReason DrawReason
}

func EvaluateGameOutcome(state *GameState) Outcome {
	legalMoves := GenerateLegalMoves(state)

	if len(legalMoves) == 0 {
		if state.IsKingInCheck() {
			if state.SideToMove == ColorWhite {
				return Outcome{Result: GameBlackWins}
			}
			return Outcome{Result: GameWhiteWins}
		}
		return Outcome{
			Result:     GameDraw,
			DrawReason: DrawStalemate,
		}
	}

	if state.HalfMoveClock >= 100 {
		return Outcome{
			Result:     GameDraw,
			DrawReason: DrawFiftyMoveRule,
		}
	}

	if hasInsufficientMaterial(state) {
		return Outcome{
			Result:     GameDraw,
			DrawReason: DrawInsufficientMaterial,
		}
	}

	//if state.IsThreefoldRepetition() {
	//	return Outcome{
	//		Result:     GameDraw,
	//		DrawReason: DrawThreefoldRepetition,
	//	}
	//}

	return Outcome{Result: GameOngoing}
}

func hasInsufficientMaterial(state *GameState) bool {
	board := state.Board

	whiteMinor := 0
	blackMinor := 0

	for sq := Square(0); sq < 64; sq++ {
		p := board[sq]
		if p.IsEmpty() {
			continue
		}

		switch p.PieceType {
		case Pawn, Queen, Rook:
			return false
		case Bishop, Knight:
			if p.Color == ColorWhite {
				whiteMinor++
			} else {
				blackMinor++
			}
		}
	}

	if whiteMinor == 0 && blackMinor == 0 {
		return true
	}

	if (whiteMinor == 1 && blackMinor == 0) ||
		(whiteMinor == 0 && blackMinor == 1) {
		return true
	}

	return false
}
