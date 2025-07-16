package enums

type GameStatus int

const (
	GameStatusOngoing GameStatus = iota
	GameStatusMafiaWin
	GameStatusPeacefulWin
)
