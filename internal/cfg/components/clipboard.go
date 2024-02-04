package components

import (
	_ "embed"
	"fmt"

	g "github.com/zaptross/gorgeous"
)

//go:embed clipboard.js
var clipboard string

func clipboardCopy(text g.JavaScript) g.JavaScript {
	g.Service("clipboard", g.JavaScript(clipboard))

	return g.JavaScript(fmt.Sprintf("copyContent(%s)", text))
}
