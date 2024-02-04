package components

import (
	_ "embed"
	"fmt"
	"regexp"

	"github.com/zaptross/countula/internal/rules"
	g "github.com/zaptross/gorgeous"
)

//go:embed rules.js
var rulesJS string

func RuleDisplay(rule rules.Rule, count int) *g.HTMLElement {
	g.Service("rules", g.JavaScript(rulesJS))

	numberClass := "rule-number"
	g.Class(&g.CSSClass{
		Selector: "." + numberClass,
		Props: g.CSSProps{
			"margin":     "0.25rem",
			"display":    "inline-block",
			"width":      "3em",
			"text-align": "center",
			"background": "var(--bg-color)",
			"color":      "var(--text-color)",
			"border":     "none",
		},
	})

	rangeClass := "rule-range"
	g.Class(&g.CSSClass{
		Selector: "." + rangeClass,
		Props: g.CSSProps{
			"margin": "0.25rem",
		},
	})

	containerClass := "rule-container"
	g.Class(&g.CSSClass{
		Selector: "." + containerClass,
		Props: g.CSSProps{
			"display":     "flex",
			"align-items": "center",
		},
	})

	g.Class(&g.CSSClass{
		Selector: "input[type=range]::-webkit-slider-runnable-track",
		Include:  true,
		Props: g.CSSProps{
			"background-color": "var(--bg-accent)",
			"border":           "1px solid var(--bg-accent)",
			"border-radius":    "5px",
		},
	})

	ruleName := rule.Name()
	ruleDesc := replaceMarkdownBold(rule.Description())

	if rule.Id() == rules.NoValidateRuleId {
		ruleName = "No Extra Validation"
		ruleDesc = "All guesses are accepted."
	}

	return g.Div(g.EB{
		Children: g.CE{
			g.H3(g.EB{Text: ruleName}),
			g.P(g.EB{Text: ruleDesc}),
			g.Div(g.EB{
				ClassList: []string{containerClass},
				Children: g.CE{
					g.Input(g.EB{
						Id:        fmt.Sprintf("rule-%d-display", rule.Id()),
						ClassList: []string{numberClass},
						Props: g.Props{
							"data-rule-id":   fmt.Sprintf("%d", rule.Id()),
							"data-rule-type": rule.Type(),
							"type":           "number",
							"min":            "0",
							"max":            "100",
							"step":           fmt.Sprintf("%d", count-1),
							"value":          fmt.Sprintf("%d", rule.Weight()),
						},
					}),
					g.Input(g.EB{
						Id:        fmt.Sprintf("rule-%d-range", rule.Id()),
						ClassList: []string{rangeClass},
						Props: g.Props{
							"data-rule-id":   fmt.Sprintf("%d", rule.Id()),
							"data-rule-type": rule.Type(),
							"type":           "range",
							"min":            "0",
							"max":            "100",
							"step":           fmt.Sprintf("%d", count-1),
							"value":          fmt.Sprintf("%d", rule.Weight()),
						},
					}),
				},
			}),
		},
	})
}

func replaceMarkdownBold(text string) string {
	return regexp.MustCompile(`\*\*([^*]+)\*\*`).ReplaceAllString(text, "<strong>$1</strong>")
}
