package components

import g "github.com/zaptross/gorgeous"

func Section(title string, children g.CE) *g.HTMLElement {
	sectionClass := "section"
	g.Class(&g.CSSClass{
		Selector: "." + sectionClass,
		Props: g.CSSProps{
			"margin": "1rem",
		},
	})
	g.Class(&g.CSSClass{
		Selector: "details > *:not(summary)",
		Include:  true,
		Props: g.CSSProps{
			"background-color": "var(--bg-color)",
			"border":           "5px solid var(--bg-color)",
		},
	})
	g.Class(&g.CSSClass{
		Selector: "details > * > h3",
		Include:  true,
		Props: g.CSSProps{
			"margin-top":  "0",
			"padding-top": "0.25rem",
		},
	})
	return g.Details(g.EB{
		ClassList: []string{sectionClass},
		Children: append(
			g.CE{
				g.Summary(g.EB{Text: title}),
			},
			children...,
		),
	})
}
