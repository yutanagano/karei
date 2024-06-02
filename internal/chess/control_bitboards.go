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

func getKingControlBitBoard(coord coordinate) bitBoard {
	return kingControlFrom[coord]
}

func getQueenControlBitBoard(coord coordinate, occupiedSquares bitBoard) bitBoard {
	controlBitBoard := bitBoard(0)

	for _, o := range []gridOffset{
		{1, 0},
		{1, 1},
		{0, 1},
		{-1, 1},
		{-1, 0},
		{-1, -1},
		{0, -1},
		{1, -1},
	} {
		controlBitBoard |= surveySlidingControlFromCoordinate(coord, o, occupiedSquares)
	}

	return controlBitBoard
}

func getRookControlBitBoard(coord coordinate, occupiedSquares bitBoard) bitBoard {
	controlBitBoard := bitBoard(0)

	for _, o := range []gridOffset{
		{1, 0},
		{0, 1},
		{-1, 0},
		{0, -1},
	} {
		controlBitBoard |= surveySlidingControlFromCoordinate(coord, o, occupiedSquares)
	}

	return controlBitBoard
}

func getBishopControlBitBoard(coord coordinate, occupiedSquares bitBoard) bitBoard {
	controlBitBoard := bitBoard(0)

	for _, o := range []gridOffset{
		{1, 1},
		{-1, 1},
		{-1, -1},
		{1, -1},
	} {
		controlBitBoard |= surveySlidingControlFromCoordinate(coord, o, occupiedSquares)
	}

	return controlBitBoard
}

func surveySlidingControlFromCoordinate(coord coordinate, offset gridOffset, occupiedSquares bitBoard) bitBoard {
	controlBitBoard := bitBoard(0)

	for {
		c, err := coord.move(offset)
		if err != nil {
			break
		}
		controlBitBoard.turnOn(c)
		if occupiedSquares.get(c) {
			break
		}
	}

	return controlBitBoard
}

func getKnightControlBitBoard(coord coordinate) bitBoard {
	return knightControlFrom[coord]
}

func getPawnKingSideControlBitBoard(pawns bitBoard, col colour) bitBoard {
	if col == white {
		return (pawns & ^bitBoardFileHMask) >> 9
	}
	return (pawns & ^bitBoardFileHMask) << 7
}

func getPawnQueenSideControlBitBoard(pawns bitBoard, col colour) bitBoard {
	if col == white {
		return (pawns & ^bitBoardFileAMask) >> 7
	}
	return (pawns & ^bitBoardFileAMask) << 9
}
