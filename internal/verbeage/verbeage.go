package verbeage

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"

	"github.com/zaptross/countula/internal/utils"
)

type Verbeage struct {
	Fail        []Response `json:"fail"`
	Count       []Response `json:"count"`
	Awaken      []Response `json:"awaken"`
	Rules       []Response `json:"rules"`
	Help        []Response `json:"help"`
	OnConfigure []Response `json:"on_configure"`
	KeepyUppies []Response `json:"keepy_uppies"`
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
	Fail        []ResponseParts
	Count       []ResponseParts
	Awaken      []ResponseParts
	Rules       []ResponseParts
	Help        []ResponseParts
	OnConfigure []ResponseParts
	KeepyUppies []ResponseParts
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
		Fail:        []ResponseParts{},
		Count:       []ResponseParts{},
		Awaken:      []ResponseParts{},
		Rules:       []ResponseParts{},
		Help:        []ResponseParts{},
		OnConfigure: []ResponseParts{},
		KeepyUppies: []ResponseParts{},
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
	for _, h := range v.Help {
		verbeage.Help = append(verbeage.Help, createResponseParts(h))
	}
	for _, o := range v.OnConfigure {
		verbeage.OnConfigure = append(verbeage.OnConfigure, createResponseParts(o))
	}
	for _, k := range v.KeepyUppies {
		verbeage.KeepyUppies = append(verbeage.KeepyUppies, createResponseParts(k))
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

func GetRandomHelpMessage() ResponseParts {
	loadTemplates()
	return utils.RandFrom(verbeage.Help)
}

func GetRandomOnConfigureMessage() ResponseParts {
	loadTemplates()
	return utils.RandFrom(verbeage.OnConfigure)
}

func GetRandomKeepyUppiesMessage() ResponseParts {
	loadTemplates()
	return utils.RandFrom(verbeage.KeepyUppies)
}
