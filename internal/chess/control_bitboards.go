package chess

var kingControlFrom [64]bitBoard
var knightControlFrom [64]bitBoard

func initKingControlBitBoards() {
	for currentSquare := coordinate(0); currentSquare < 64; currentSquare++ {
		controlBitBoard := bitBoard(0)

		for _, d := range []gridOffset{
			{1, 0},
			{1, 1},
			{0, 1},
			{-1, 1},
			{-1, 0},
			{-1, -1},
			{0, -1},
			{1, -1},
		} {
			if controlledSquare, err := currentSquare.move(d); err == nil {
				controlBitBoard.turnOn(controlledSquare)
			}
		}

		kingControlFrom[currentSquare] = controlBitBoard
	}
}

func initKnightControlBitBoards() {
	for currentSquare := coordinate(0); currentSquare < 64; currentSquare++ {
		controlBitBoard := bitBoard(0)

		for _, d := range []gridOffset{
			{2, -1},
			{2, 1},
			{1, 2},
			{-1, 2},
			{-2, -1},
			{-2, 1},
			{1, -2},
			{-1, -2},
		} {
			if controlledSquare, err := currentSquare.move(d); err == nil {
				controlBitBoard.turnOn(controlledSquare)
			}
		}

		knightControlFrom[currentSquare] = controlBitBoard
	}
}
