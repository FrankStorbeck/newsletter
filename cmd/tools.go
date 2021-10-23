package main

import (
	"html/template"
	"io/ioutil"
)

// getTemplate parses the contents of a file holding the template and returns
// the result
func getTemplate(file, name string) (*template.Template, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return template.New(name).Parse(string(b))
}
