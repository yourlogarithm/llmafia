package enums

type Phase string

const (
	PhaseDayDiscussion               Phase = "Day Discussion"
	PhaseDayVoting                   Phase = "Day Voting"
	PhaseNightDoctorHealing          Phase = "Night Doctor Healing"
	PhaseNightDetectiveInvestigation Phase = "Night Detective Investigation"
	PhaseNightMafiaDiscussion        Phase = "Night Mafia Discussion"
	PhaseNightMafiaElimination       Phase = "Night Mafia Elimination"
)
