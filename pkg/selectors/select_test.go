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
		{"Missing", "", ErrNoMatch},
		{"OK", "1;John;O';Doe;john@company.com;Street;10;8900;un;Yes;comment", nil},
		{"Wrong ID", "2_;John;O';Doe;john@company.com;;;;;Yes", ErrNoMatch},
		{"Wrong family name", "3;John;O';Did;john@company.com;;;;;Yes", ErrNoMatch},
		{"No email", "4;John;O';Doe;;Street;10;8900;un;Yes", ErrNoMatch},
		{"Wrong email", "5;John;O';Doe;john;Street;10;8900;un;Yes", ErrInvalidEMail},
	}

	colNames := []string{"id", "first name", "middle names", "family name", "email", "street", "number", "zip", "country", "wants newsletter"}

	for _, tst := range tests {
		record := strings.Split(tst.line, ";")
		_, err := slctrs.Select(record, colNames)
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
		}
	}
}
