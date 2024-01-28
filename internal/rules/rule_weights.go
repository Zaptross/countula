package rules

type RuleWeight struct {
	Current int
	Default int
}

func Weights(weight int) RuleWeight {
	return RuleWeight{
		Current: weight,
		Default: weight,
	}
}
