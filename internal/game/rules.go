package game

import "github.com/zaptross/countula/internal/rules"

func getRulesForNewGame() int {
	pv := rules.GetRandomPreValidateRule()
	v := rules.GetRandomValidateRule()

	return pv.Id() | v.Id()
}
