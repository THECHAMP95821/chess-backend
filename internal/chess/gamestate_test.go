package chess

import (
	"fmt"
	"testing"
)

func TestInitialGameState(t *testing.T) {
	initialGameState := NewInitialGameState()
	fmt.Println(NewSquare(4, 3).String())
	initialGameState.Board.DisplayBoard()
}
