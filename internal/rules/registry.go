package rules

import "github.com/zaptross/countula/internal/emoji"

const (
	// Order matters here
	IncrementRule1Id     = 1 << iota
	IncrementRule2Id     = 1 << iota
	IncrementRule3Id     = 1 << iota
	IncrementRule7Id     = 1 << iota
	TakeTurnsRuleId      = 1 << iota
	GuessNormallyRuleId  = 1 << iota
	NoValidateRuleId     = 1 << iota
	FibonacciRuleId      = 1 << iota
	RomanNumeralRuleId   = 1 << iota
	JeopardyRuleId       = 1 << iota
	GoodyTwoShoesRuleId  = 1 << iota
	KeepyUppiesRuleId    = 1 << iota
	CountOfTheHillRuleId = 1 << iota
)

// Default Rule Weights
const (
	// Pre-Validate Rules
	GuessNormallyRuleWeight = 56
	RomanNumeralRuleWeight  = 28
	JeopardyRuleWeight      = 16

	// Count Rules
	IncrementRule1Weight     = 16
	IncrementRule2Weight     = 20
	IncrementRule3Weight     = 20
	IncrementRule7Weight     = 12
	FibonacciRuleWeight      = 12
	GoodyTwoShoesRuleWeight  = 10
	CountOfTheHillRuleWeight = 10

	// Validate Rules
	NoValidateRuleWeight  = 30
	TakeTurnsRuleWeight   = 55
	KeepyUppiesRuleWeight = 15
)

func OverrideSuccessEmoji(ruleId int) string {
	switch ruleId {
	case KeepyUppiesRuleId:
		return emoji.BALLOON
	}
	return emoji.CHECK
}

func OverrideRuleSelections(rules int) int {
	switch true {
	case rules&(TakeTurnsRuleId|CountOfTheHillRuleId) == TakeTurnsRuleId|CountOfTheHillRuleId:
		// Replace Take Turns with NoValidate while Count of the Hill is active,
		// because Count of the Hill has custom behavior for taking turns.
		return rules ^ (TakeTurnsRuleId | NoValidateRuleId)
	}
	return rules
}
