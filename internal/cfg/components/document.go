package components

import (
	g "github.com/zaptross/gorgeous"
)

func Head() *g.HTMLElement {
	return g.Head(g.EB{
		Children: g.CE{},
	})
}

func Body() *g.HTMLElement {
	return g.Body(g.EB{
		Children: g.CE{
			g.H1(g.EB{Text: "Yeahnahmate"}),
		},
	})
}
