package main

import (
	"os"

	g "github.com/zaptross/gorgeous"
)

func createDistDirectories() {
	os.Mkdir("dist", 0755)
	os.Mkdir("dist/public", 0755)
}

func writeRenderedHTML(rendered *g.RenderedHTML) {
	os.WriteFile("dist/index.html", []byte(rendered.Document), 0644)
	os.WriteFile("dist/style.css", []byte(rendered.Style), 0644)
	os.WriteFile("dist/script.js", []byte(rendered.Script), 0644)
}

func copyPublicToDist() {
	files, err := os.ReadDir("./public")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileContents, fcErr := os.ReadFile("public/" + file.Name())

		if fcErr != nil {
			panic(fcErr)
		}

		os.WriteFile("dist/public/"+file.Name(), fileContents, 0644)
	}
}
