package chess

type MoveFlags uint8

const (
	MoveFlagNone       = MoveFlags(0)
	MoveFlagCapture    = MoveFlags(1 << 0)
	MoveFlagEnPassant  = MoveFlags(1 << 1)
	MoveFlagCastle     = MoveFlags(1 << 2)
	MoveFlagPromotion  = MoveFlags(1 << 3)
	MoveFlagDoublePush = MoveFlags(1 << 4)
)

type Move struct {
	From      Square
	To        Square
	Flags     MoveFlags
	Promotion PieceType
}

func NewMove(from, to Square) *Move {
	return &Move{
		From:  from,
		To:    to,
		Flags: MoveFlagNone,
	}
}
func NewCaptureMove(from, to Square) *Move {
	return &Move{
		From:  from,
		To:    to,
		Flags: MoveFlagCapture,
	}
}
func NewPromotionMove(from, to Square, promoPiece PieceType) *Move {
	return &Move{
		From:      from,
		To:        to,
		Flags:     MoveFlagPromotion,
		Promotion: promoPiece,
	}
}
func NewPromotionCapture(from, to Square) *Move {
	return &Move{
		From:  from,
		To:    to,
		Flags: MoveFlagCapture | MoveFlagEnPassant,
	}
}
func NewEnPassant(from, to Square) *Move {
	return &Move{
		From:  from,
		To:    to,
		Flags: MoveFlagEnPassant | MoveFlagCapture,
	}
}
func NewCastleMove(from, to Square) *Move {
	return &Move{
		From:  from,
		To:    to,
		Flags: MoveFlagCastle,
	}
}
func NewDoublePush(from, to Square) *Move {
	return &Move{
		From:  from,
		To:    to,
		Flags: MoveFlagDoublePush,
	}
}
func (m Move) IsCapture() bool {
	return m.Flags&MoveFlagCapture != 0
}
func (m Move) IsEnPassant() bool {
	return m.Flags&MoveFlagEnPassant != 0
}
func (m Move) IsCastle() bool {
	return m.Flags&MoveFlagCastle != 0
}
func (m Move) IsPromotion() bool {
	return m.Flags&MoveFlagPromotion != 0
}
func (m Move) IsDoublePush() bool {
	return m.Flags&MoveFlagDoublePush != 0
}
func (m Move) String() string {
	move := m.From.String()
	move += m.To.String()
	if m.Promotion != PieceNone {
		move += m.Promotion.String()
	}
	return move
}
func ParseMove(s string) (Move, string) {
	if len(s) < 4 {
		return Move{}, "invalid move string"
	}

	from := ParseSquare(s[0:2])
	to := ParseSquare(s[2:4])

	move := Move{
		From: from,
		To:   to,
	}

	if len(s) == 5 {
		move.Promotion = ParsePieceType(string(s[4]))
	}

	return move, ""
}
