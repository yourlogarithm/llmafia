package prompts

var common_prompt = `You are a member of a Mafia social deduction game. The game consists of 8 players:
* 2 Mafia members
* 1 Doctor
* 1 Detective
* 4 Citizens

There is also a narrator who describes the game state.
Your name is %s and your role is %s.`

var caution_note = "Note that your response is entirely visible for other players during day phases, so be careful not to include revealing information (ex. inner thoughts) in your responses."

var peaceful_vote = "You can propose people for vote, if you think they are Mafia. If the majority of players agree, that person will be eliminated from the game."
