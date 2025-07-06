package enums

type Role string

const (
	RoleCitizen   Role = "Citizen"
	RoleDetective Role = "Detective"
	RoleDoctor    Role = "Doctor"
	RoleMafia     Role = "Mafia"
	RoleNarrator  Role = "Narrator" // Narrator is the one who manages the game flow and announces phases
)
