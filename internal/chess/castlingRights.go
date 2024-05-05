package chess

type castlingRights uint8

const (
	whiteCastleKingSide  castlingRights = 0b0001
	whiteCastleQueenSide castlingRights = 0b0010
	blackCastleKingSide  castlingRights = 0b0100
	blackCastleQueenSide castlingRights = 0b1000
)

func (c castlingRights) isSet(flag castlingRights) bool {
	return c&flag != 0
}

func (c *castlingRights) turnOn(flag castlingRights) {
	*c |= flag
}

func (c *castlingRights) turnOff(flag castlingRights) {
	*c &= ^flag
}
