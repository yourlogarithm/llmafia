package enums

type GameStatus int

const (
	GameStatusMafiaWin GameStatus = iota
	GameStatusPeacefulWin
	GameStatusOngoing
)
