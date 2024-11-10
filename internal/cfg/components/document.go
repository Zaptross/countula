package components

import (
	"sort"

	"github.com/samber/lo"
	"github.com/zaptross/countula/internal/rules"
	g "github.com/zaptross/gorgeous"
)

func Head() *g.HTMLElement {
	return g.Head(g.EB{
		Children: g.CE{
			g.Link(g.EB{
				Props: g.Props{
					"rel":  "stylesheet",
					"type": "text/css",
					"href": "style.css",
				},
			}),
			g.Script(g.EB{
				Props: g.Props{
					"type": "text/javascript",
					"src":  "script.js",
				},
			}),
		},
	})
}

func Body() *g.HTMLElement {
	g.Class(&g.CSSClass{
		Selector: ":root",
		Include:  true,
		Props: g.CSSProps{
			"--bg-color":   BackgroundColour,
			"--bg-accent":  BackgroundAccent,
			"--text-color": FontColour,
		},
	})
	g.Class(&g.CSSClass{
		Selector: "html, body",
		Include:  true,
		Props: g.CSSProps{
			"color":       "var(--text-color)",
			"font-family": "sans-serif",
			"font-size":   "20px",
		},
	})
	g.Class(&g.CSSClass{
		Selector: "html",
		Include:  true,
		Props: g.CSSProps{
			"background-color": "var(--bg-color)",
		},
	})
	g.Class(&g.CSSClass{
		Selector: "body",
		Include:  true,
		Props: g.CSSProps{
			"background-color": "var(--bg-accent)",
			"margin":           "1rem auto 0 auto",
			"width":            "max(50%, min(100%, calc(40vw + 270px)))",
			"border":           "5px solid var(--bg-accent)",
			"border-radius":    "5px",
		},
	})

	countingRules := rules.FilterForType(rules.CountType)(rules.AllRules)
	preValidateRules := rules.FilterForType(rules.PreValidateType)(rules.AllRules)
	validateRules := rules.FilterForType(rules.ValidateType)(rules.AllRules)

	sortById := func(i, j int) bool {
		return countingRules[i].Id() < countingRules[j].Id()
	}

	sort.Slice(countingRules, sortById)
	sort.Slice(preValidateRules, sortById)
	sort.Slice(validateRules, sortById)

	preValidateDiv := Section(
		"Pre-Validate Rules",
		"Pre-Validate rules are checked before the number is incremented and they change the way the user submits the number.",
		lo.Map(preValidateRules, func(rule rules.Rule, _ int) *g.HTMLElement {
			return RuleDisplay(rule, len(preValidateRules))
		}),
	)

	countingDiv := Section(
		"Counting Rules",
		"Counting rules change the way the next number is decided.",
		lo.Map(countingRules, func(rule rules.Rule, _ int) *g.HTMLElement {
			return RuleDisplay(rule, len(countingRules))
		}),
	)

	validateDiv := Section(
		"Validate Rules",
		"Validate rules are checked once the user's guess is confirmed to match the next number, and they change meta aspects of the game.",
		lo.Map(validateRules, func(rule rules.Rule, _ int) *g.HTMLElement {
			return RuleDisplay(rule, len(validateRules))
		}),
	)

	rulesDisplayClass := "rules-display"
	g.Class(&g.CSSClass{
		Selector: "." + rulesDisplayClass + " > details > summary",
		Include:  true,
		Props: g.CSSProps{
			"cursor":      "pointer",
			"font-size":   "1.5rem",
			"font-weight": "bold",
		},
	})

	rulesDisplay := g.Div(g.EB{
		ClassList: []string{rulesDisplayClass},
		Children: g.CE{
			preValidateDiv,
			countingDiv,
			validateDiv,
		},
	})

	return g.Body(g.EB{
		Children: g.CE{
			InstructionsAndTitle(),
			CommandDisplay(),
			rulesDisplay,
		},
	})
}
