package selectors

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

const (
	// EMail is a name of the field holding a valid e-mail address for a recipient
	EMail = "EMail"
)

var (
	rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+" +
		"@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9]" +
		"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Recipient maps the field names to their value
type Recipient map[string]string

// Get returns the value for the key.
func (rcpnt Recipient) Get(key string) string {
	return rcpnt[key]
}

// Set sets the value for key
func (rcpnt Recipient) Set(key, value string) {
	rcpnt[strings.TrimSpace(key)] = strings.TrimSpace(value)
}

// selector is a struct for testing a value in a field
type selector struct {
	fieldName string // name of the field
	value     string // value that the field must have for selection
}

// Selectors holds the Selectors for fields
type Selectors []selector

// Read reads the Selectors from a file located at path. Each line
// holds the name of the field, then a `=` charcter and then the value a field
// must have in order to be selected as a valid recipient. The order of the
// lines should be the same as the order of the fields (collumns) in the the
// files with the recipents.
func Read(path string) (*Selectors, error) {
	slctrs := Selectors{}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	i := 0
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return nil, err
		}

		line := scanner.Text()
		fields := strings.Split(line, "=")
		if len(fields) >= 2 {
			slctrs = append(slctrs, selector{fields[0], fields[1]})
		}
		i++
	}

	return &slctrs, nil
}

// TestRecipient tests if line holds a recipient eligible for sending a
// newsletter. The test passes when all Selectors find a correct field value and
// when the recipient has a valid e-mail address.
func (slctrs Selectors) TestRecipient(line string) (*Recipient, bool) {
	fields := strings.Split(line, ";")
	rcpnt := Recipient(make(map[string]string, 0))

	ls := len(slctrs)
	lf := len(fields)
	for i := 0; i < ls; i++ {
		if i >= lf {
			// missing fields never match
			return nil, false
		}
		v := slctrs[i].value
		if v != "*" && fields[i] != v {
			return nil, false
		}
		rcpnt.Set(slctrs[i].fieldName, fields[i])
	}
	if em := rcpnt.Get(EMail); len(em) == 0 || !rxEmail.MatchString(em) {
		return nil, false
	}
	// missing selectors always match
	return &rcpnt, true
}
