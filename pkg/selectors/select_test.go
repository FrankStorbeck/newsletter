package selectors

import (
	"strings"
	"testing"
)

var path = "../../tests/selectors.txt"

func TestSelect(t *testing.T) {
	slctrs, err := New(path)
	if err != nil {
		t.Fatalf("For Select: New(%q) returns an error: %s", path, err.Error())
	}

	tests := []struct {
		name string
		line string
		err  error
	}{
		{"OK", "1;John;O';Doe;john@company.com;Street;10;8900;un;Yes;comment", nil},
		{"Wrong family name", "3;John;O';Did;john@company.com;;;;;Yes", ErrNoMatch},
		{"No e-mail", "4;John;O';Doe;;Street;10;8900;un;Yes", ErrNoMatch},
		{"Wrong e-mail", "5;John;O';Doe;john;Street;10;8900;un;Yes", ErrInvalidEMail},
	}

	colNames := []string{"id", "first name", "middle names", "family name",
		"email", "street", "number", "zip", "country", "wants newsletter"}

	for _, tst := range tests {
		record := strings.Split(tst.line, ";")
		recipient, err := slctrs.Select(record, colNames, 4)
		switch {
		case tst.err == nil && err != nil:
			t.Errorf("%s: Select() returns an error %q, should be nil",
				tst.name, err.Error())
		case tst.err != nil && err == nil:
			t.Errorf("%s: Select() returns no error, should be %q",
				tst.name, tst.err.Error())
		case tst.err != nil && err != nil:
			if tst.err != err {
				t.Errorf("%s: Select() returns error %q, should be %q",
					tst.name, err.Error(),
					tst.err.Error())
			}
		default:
			for i, name := range colNames {
				if got := (*recipient).Get(name); got != record[i] {
					t.Fatalf("%s: After Select() (*recipient).Get(%q) is %q, should be %q",
						tst.name, name, got, record[i])
				}
			}
		}
	}
}
