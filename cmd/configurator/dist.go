package main

import (
	"os"

	g "github.com/zaptross/gorgeous"
)

func createDistDirectories() {
	os.Mkdir("dist", 0755)
}

func writeRenderedHTML(rendered *g.RenderedHTML) {
	os.WriteFile("dist/index.html", []byte(rendered.Document), 0644)
	os.WriteFile("dist/style.css", []byte(rendered.Style), 0644)
	os.WriteFile("dist/script.js", []byte(rendered.Script), 0644)
}
