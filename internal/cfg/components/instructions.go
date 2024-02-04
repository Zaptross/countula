package components

import (
	g "github.com/zaptross/gorgeous"
)

func InstructionsAndTitle() *g.HTMLElement {
	titleClass := "page-title"
	g.Class(&g.CSSClass{
		Selector: "." + titleClass,
		Props: g.CSSProps{
			"margin": "0.25rem auto",
			"width":  "fit-content",
		},
	})

	instructionsClass := "instructions"
	g.Class(&g.CSSClass{
		Selector: "." + instructionsClass,
		Props: g.CSSProps{
			"padding": "0.5rem",
		},
	})

	return g.Div(g.EB{
		Children: g.CE{
			g.H1(g.EB{
				ClassList: []string{titleClass},
				Text:      "Countula Configurator",
			}),
			g.Div(g.EB{
				ClassList: []string{instructionsClass},
				Children: g.CE{
					g.P(g.EB{Text: "Adjust the weights of the rules below to change how how frequent each rule is."}),
					g.P(g.EB{Text: "To apply these settings, copy the command below and paste it into the channel where Countula is running."}),
				},
			}),
		},
	})
}
