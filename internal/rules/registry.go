package rules

const (
	// Order matters here
	IncrementRule1Id    = 1 << iota
	IncrementRule2Id    = 1 << iota
	IncrementRule3Id    = 1 << iota
	IncrementRule7Id    = 1 << iota
	TakeTurnsRuleId     = 1 << iota
	GuessNormallyRuleId = 1 << iota
	NoValidateRuleId    = 1 << iota
	FibonacciRuleId     = 1 << iota
	RomanNumeralRuleId  = 1 << iota
	JeopardyRuleId      = 1 << iota
	GoodyTwoShoesRuleId = 1 << iota
)

// Default Rule Weights
const (
	// Pre-Validate Rules
	GuessNormallyRuleWeight = 56
	RomanNumeralRuleWeight  = 28
	JeopardyRuleWeight      = 16

	// Count Rules
	IncrementRule1Weight    = 22
	IncrementRule2Weight    = 22
	IncrementRule3Weight    = 18
	IncrementRule7Weight    = 14
	FibonacciRuleWeight     = 14
	GoodyTwoShoesRuleWeight = 10

	// Validate Rules
	NoValidateRuleWeight = 40
	TakeTurnsRuleWeight  = 60
)
