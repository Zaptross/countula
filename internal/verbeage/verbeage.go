package verbeage

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"
	"math/rand"
)

type Verbeage struct {
	Fail  []Response `json:"fail"`
	Count []Response `json:"count"`
}

type Response struct {
	Reply   *string `json:"reply"`
	Message *string `json:"message,omitempty"`
}

type ExecutableTemplate func(TemplateFields) (string, error)
type ResponseParts struct {
	Reply   ExecutableTemplate
	Message ExecutableTemplate
}

type ResponseTemplate struct {
	Fail  []ResponseParts
	Count []ResponseParts
}

type TemplateFields struct {
	Username string
}

//go:embed verbeage.json
var rawVerbeage string

var verbeage ResponseTemplate

func getVerbeage() Verbeage {
	var v Verbeage
	json.Unmarshal([]byte(rawVerbeage), &v)
	return v
}

func loadTemplates() {
	if verbeage.Fail != nil {
		return
	}

	v := getVerbeage()
	verbeage = ResponseTemplate{
		Fail:  []ResponseParts{},
		Count: []ResponseParts{},
	}

	for _, f := range v.Fail {
		verbeage.Fail = append(verbeage.Fail, createResponseParts(f))
	}
	for _, c := range v.Count {
		verbeage.Count = append(verbeage.Count, createResponseParts(c))
	}
}

func createExecutableTemplate(t string) ExecutableTemplate {
	executableTemplate := template.Must(template.New("").Parse(t))
	return func(tf TemplateFields) (string, error) {
		var buf bytes.Buffer
		err := executableTemplate.Execute(&buf, tf)
		return buf.String(), err
	}
}

func createResponseParts(r Response) ResponseParts {
	rp := ResponseParts{}

	if r.Reply != nil {
		rp.Reply = createExecutableTemplate(*r.Reply)
	}

	if r.Message != nil {
		rp.Message = createExecutableTemplate(*r.Message)
	}

	return rp
}

func randFrom[T any](arr []T) T {
	return arr[rand.Intn(len(arr))]
}

func GetRandomFail() ResponseParts {
	loadTemplates()
	return randFrom(verbeage.Fail)
}
