package chess

var kingControlFrom [64]bitBoard
var knightControlFrom [64]bitBoard

func initKingControlBitBoards() {
	for currentSquare := a1; currentSquare <= h8; currentSquare++ {
		controlBitBoard := bitBoard(0)

		for _, o := range []offset{
			{1, 0},
			{1, 1},
			{0, 1},
			{-1, 1},
			{-1, 0},
			{-1, -1},
			{0, -1},
			{1, -1},
		} {
			if controlledSquare, err := currentSquare.move(o); err == nil {
				controlBitBoard.turnOn(controlledSquare)
			}
		}

		kingControlFrom[currentSquare] = controlBitBoard
	}
}

func initKnightControlBitBoards() {
	for currentSquare := a1; currentSquare <= h8; currentSquare++ {
		controlBitBoard := bitBoard(0)

		for _, o := range []offset{
			{2, -1},
			{2, 1},
			{1, 2},
			{-1, 2},
			{-2, -1},
			{-2, 1},
			{1, -2},
			{-1, -2},
		} {
			if controlledSquare, err := currentSquare.move(o); err == nil {
				controlBitBoard.turnOn(controlledSquare)
			}
		}

		knightControlFrom[currentSquare] = controlBitBoard
	}
}
