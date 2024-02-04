package components

import (
	"fmt"

	g "github.com/zaptross/gorgeous"
)

func CommandDisplay() *g.HTMLElement {
	commandClass := "command-copy"
	g.Class(&g.CSSClass{
		Selector: "." + commandClass,
		Props: g.CSSProps{
			"cursor":           "pointer",
			"margin":           "0.25rem",
			"padding":          "0.75rem",
			"border":           "5px solid var(--bg-color)",
			"border-radius":    "5px",
			"background-color": "var(--bg-color)",
			"font-family":      "monospace",
		},
	})

	return g.P(g.EB{
		Id:        "rule-command",
		ClassList: []string{commandClass},
		Text:      "/count settings 11:11,11:11,11:11,11:11,11:11,11:11,11:11,11:11",
		Script: g.JavaScript(fmt.Sprintf(
			`thisElement.onclick = ()=>{%s};`,
			clipboardCopy(g.JavaScript(`thisElement.innerText`)),
		)),
	})
}
