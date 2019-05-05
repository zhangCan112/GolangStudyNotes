package negroni

import (
	"os"
	"testing"
	"text/template"
)

func TestTempmate(t *testing.T) {
	const letter = `Dear {{.Name}},

{{if .Attended}}It was a pleasure to see you at the wedding.
如果Attended是true的话，这句是第二行{{else}}It is a shame you couldn't make it to the wedding.
如果Attended是false的话，这句是第二行{{end}}
{{with .Gift}}Thank you for the lovely {{.}}.
{{end}}
Best wishes,
Josie

`
	tmpl, err := template.New("letter").Parse(letter)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, nil)
	if err != nil {
		panic(err)
	}
}
