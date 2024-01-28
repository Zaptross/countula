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
)

// Default Rule Weights
const (
	// Count Rules
	IncrementRule1Weight = 25
	IncrementRule2Weight = 25
	IncrementRule3Weight = 20
	IncrementRule7Weight = 15
	FibonacciRuleWeight  = 15

	// Validate Rules
	NoValidateRuleWeight   = 55
	RomanNumeralRuleWeight = 27
	JeopardyRuleWeight     = 18

	// Pre-Validate Rules
	TakeTurnsRuleWeight     = 50
	GuessNormallyRuleWeight = 50
)
