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

	animDuration := 0.5
	g.Class(&g.CSSClass{
		Selector: ".copied-anim",
		Include:  true,
		Props: g.CSSProps{
			"animation": fmt.Sprintf("%fs ease-in-out 1 copied", animDuration),
		},
	})

	g.Class(&g.CSSClass{
		Include: true,
		Raw: `
		@keyframes copied {
			0% {
				border-color: rgb(70, 75, 70);
				background-color: rgb(50, 55, 50);
			}
			100% {}
		}
		`,
	})

	return g.P(g.EB{
		Id:        "rule-command",
		ClassList: []string{commandClass},
		Text:      "/count settings action:set settings:11:11,11:11,11:11,11:11,11:11,11:11,11:11,11:11",
		Script: g.JavaScript(fmt.Sprintf(
			`thisElement.onclick = ()=>{
				if (!thisElement.classList.contains("copied-anim")) {
					thisElement.classList.add("copied-anim");
					setTimeout(()=>thisElement.classList.remove("copied-anim"), %f * 1000);
				}

				%s
			};`,
			animDuration,
			clipboardCopy(g.JavaScript(`thisElement.innerText`)),
		)),
	})
}
