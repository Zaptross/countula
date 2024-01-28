package main

import (
	c "github.com/zaptross/countula/internal/cfg/components"
	g "github.com/zaptross/gorgeous"
)

func main() {
	rendered := g.RenderStatic(
		g.Document(
			c.Head(),
			c.Body(),
		),
	)

	createDistDirectories()
	writeRenderedHTML(rendered)
	copyPublicToDist()
}
