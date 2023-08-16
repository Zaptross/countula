package game

import "github.com/zaptross/countula/internal/rules"

func getRulesForNewGame() int {
	pv := rules.GetRandomPreValidateRule()
	c := rules.GetRandomCountRule()
	v := rules.GetRandomValidateRule()

	return pv.Id() | c.Id() | v.Id()
}
