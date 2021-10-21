package main

import (
	"html/template"
	"strings"
)

// parseTemplates returns the parsed templates found in tmplts. The map holding
// pointers to the parsed templates have keys equal to their extension.
func parseTemplates(tmplts []string) (map[string]*template.Template, error) {
	parsed := make(map[string]*template.Template)

	for _, t := range tmplts {
		if i := strings.LastIndex(t, "."); i > 0 {
			var err error
			parsed[t[i:]], err = template.ParseFiles(t)
			if err != nil {
				return parsed, err
			}
		}
	}

	return parsed, nil
}
