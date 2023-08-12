package verbeage

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"

	"github.com/zaptross/countula/internal/utils"
)

type Verbeage struct {
	Fail   []Response `json:"fail"`
	Count  []Response `json:"count"`
	Awaken []Response `json:"awaken"`
	Rules  []Response `json:"rules"`
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
	Fail   []ResponseParts
	Count  []ResponseParts
	Awaken []ResponseParts
	Rules  []ResponseParts
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
		Fail:   []ResponseParts{},
		Count:  []ResponseParts{},
		Awaken: []ResponseParts{},
		Rules:  []ResponseParts{},
	}

	for _, f := range v.Fail {
		verbeage.Fail = append(verbeage.Fail, createResponseParts(f))
	}
	for _, c := range v.Count {
		verbeage.Count = append(verbeage.Count, createResponseParts(c))
	}
	for _, a := range v.Awaken {
		verbeage.Awaken = append(verbeage.Awaken, createResponseParts(a))
	}
	for _, r := range v.Rules {
		verbeage.Rules = append(verbeage.Rules, createResponseParts(r))
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

func GetRandomFail() ResponseParts {
	loadTemplates()
	return utils.RandFrom(verbeage.Fail)
}

func GetRandomCount() ResponseParts {
	loadTemplates()
	return utils.RandFrom(verbeage.Count)
}

func GetRandomAwaken() ResponseParts {
	loadTemplates()
	return utils.RandFrom(verbeage.Awaken)
}

func GetRandomRuleMessage() ResponseParts {
	loadTemplates()
	return utils.RandFrom(verbeage.Rules)
}
