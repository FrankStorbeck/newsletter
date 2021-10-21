package selectors

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
	"text/template"
)

func TestSelect(t *testing.T) {
	slctrs, err := Read("../../tests/selectors_test.txt")
	if err != nil {
		t.Fatalf("Read() returns an error: %s", err.Error())
	}

	f, err := os.Open("../../tests/recipients_test.csv")
	if err != nil {
		t.Fatalf("os.Open() returns an error: %s", err.Error())
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if i := strings.Index(line, "#"); i >= 0 {
			line = line[:i]
		}
		if len(line) > 0 {
			rcpnt, ok := slctrs.TestRecipient(scanner.Text())
			if ok {
				tmpl, err := template.New("test").Parse("{{.Get \"EMail\" }} belongs to {{.Get \"FirstName\"}} {{.Get \"MiddleNames\"}} {{.Get \"FamilyName\"}}")
				if err != nil {
					t.Fatalf("Cannot create template test: %s", err.Error())
				}
				body := new(bytes.Buffer)
				err = tmpl.Execute(body, rcpnt)
				if err != nil {
					t.Fatalf("Cannot create output from template: %s", err.Error())
				}
				want := "frank@storbeck.nl belongs to F.  Storbeck"
				if got := body.String(); got != want {
					t.Errorf("Result is\n%q\nshould be\n%q", got, want)
				}
			} else {
				t.Errorf("failure")
			}
		}
	}
}
