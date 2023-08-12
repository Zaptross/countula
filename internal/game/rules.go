package game

import "github.com/zaptross/countula/internal/rules"

func getRulesForNewGame() int {
	return rules.IncrementOne.Id() | rules.GuessNormally.Id()
}
