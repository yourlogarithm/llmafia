# LLMafia

A command-line toy to force your LLM to play against itself.

## Overview

LLMafia simulates a complete Mafia game with 7 AI players, each with distinct roles and objectives. The game uses OpenAI's API to generate realistic dialogue and strategic decisions, creating engaging gameplay where AI players attempt to deduce, deceive, and survive.

## Game Setup

### Players & Roles

Each game consists of **7 players** with the following role distribution:

- **2 Mafia Members** - Secret partners working to eliminate all other players
- **1 Detective** - Can investigate one player each night to learn their alignment
- **1 Doctor** - Can protect one player from elimination each night
- **3 Citizens** - Regular townspeople trying to identify and eliminate the Mafia

### Player Names

Players are randomly assigned from a pool of names: Alice, Bob, Charlie, Diana, Ethan, Fiona, and George.

## Game Cycle

The game alternates between **Day Phases** and **Night Phases** until one side achieves victory.

### Day Phase

1. **Discussion Round**: Players speak one by one in rotated order, sharing observations and suspicions
2. **Accusations**: Players may accuse others of being Mafia during discussions
3. **Voting**: If accusations were made, players vote to eliminate suspected Mafia members
   - A player needs >50% of votes to be eliminated
   - Players may abstain from voting
   - If no player receives majority votes, no one is eliminated

### Night Phase

Three actions occur simultaneously during the night:

1. **Mafia Elimination**: The two Mafia members coordinate to select a target for elimination
   - If both Mafia agree, their target is eliminated
   - If they disagree, one of their choices is randomly selected
   - If only one Mafia remains, they choose alone

2. **Doctor Protection**: The Doctor selects one player to protect from elimination
   - Can protect anyone including themselves, but not twice in consecutive nights
   - If the Doctor protects the Mafia's target, that player survives

3. **Detective Investigation**: The Detective chooses one player to investigate
   - Learns whether the target is "Mafia" or "Not Mafia"
   - Uses this information to guide future discussions and votes

## Victory Conditions

- **Mafia Wins**: When Mafia members equal or outnumber the remaining peaceful players
- **Peaceful Players Win**: When all Mafia members have been eliminated

## Game Rules & Strategy

### Role-Specific Objectives

**Mafia Members:**
- Eliminate all non-Mafia players while avoiding detection
- Coordinate night eliminations with their partner
- Blend in during day discussions and deflect suspicion
- May propose eliminations during day phase to influence voting

**Detective:**
- Gather intelligence through nightly investigations
- Use information strategically without revealing their role too early
- Guide town voting toward Mafia elimination
- Balance between staying hidden and sharing crucial information

**Doctor:**
- Protect key players from Mafia elimination
- Observe behavior to identify likely Mafia targets
- Decide when (if ever) to reveal their role for town benefit
- Cannot protect the same player on consecutive nights (except themselves)

**Citizens:**
- Analyze discussions for suspicious behavior and contradictions
- Use voting power strategically to eliminate suspected Mafia
- May claim false roles to protect actual power roles
- Rely on logical deduction and social reads

### Communication Rules

- During day phases, all responses are visible to every player
- Players must stay completely in character
- No meta-commentary or third-person descriptions allowed
- Responses should sound like natural human dialogue
- Strategic deception and misdirection are encouraged

## Running the Game

```bash
go run cmd ./cmd/ --base-url https://llm.url --model llm-model-name --out file.json
```

The `--out` flag is optional and specifies an output file for game logs, if not provided, game logs won't be saved.