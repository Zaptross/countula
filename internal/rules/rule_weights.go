package rules

type RuleWeight struct {
	Current int
}

func SetupWeight(weight int) RuleWeight {
	return RuleWeight{
		Current: weight,
	}
}
