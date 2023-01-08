package selectors

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

const (
	// EMailColName is a name of the column holding a valid e-mail address for a recipient
	EMailColName = "email"
)

// Errors
var (
	ErrFieldsMissing = errors.New("Too few fields")
	ErrInvalidEMail  = errors.New("Invalid e-mail")
	ErrNoMatch       = errors.New("No match")

	rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+" +
		"@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9]" +
		"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Selectors holds the Selectors for columns
type Selectors map[string]*regexp.Regexp

// New returns a new slice with selectors. pathToSelectorsFile is the path
// to a selectors file. Each line in it should contain a column name followed by
// an `=` character and a regular expression to be used for testing if the value
// in the column matches. The regular expression must be surrounded by quoting
// characters (`"`).
func New(pathToSelectorsFile string) (*Selectors, error) {
	f, err := os.Open(pathToSelectorsFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	slctrs := make(Selectors, 0)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return nil, err
		}

		fields := strings.Split(scanner.Text(), "=")
		if len(fields) >= 2 {
			if colName := strings.ToLower(strings.TrimSpace(fields[0])); len(colName) > 0 {
				re := ""
				if i := strings.IndexByte(fields[1], '"'); i >= 0 {
					if j := strings.LastIndexByte(fields[1], '"'); j > i {
						re = fields[1][i+1 : j]
					}
				}
				if len(re) > 0 {
					slctrs[colName], err = regexp.Compile(re)
					if err != nil {
						return &slctrs, err
					}
				}
			}
		}
	}

	return &slctrs, nil
}

// Select tests if a subscriber is eligible for receiving a newsletter. fields
// must hold the field values for the subscriber's record and colNames the
// collumn names for these fields. The length of fields will be truncated to the
// length of colNames.
func (slctrs Selectors) Select(fields, colNames []string) (*Recipient, error) {
	l := len(fields)
	if l > len(colNames) {
		l = len(colNames)
	}

	rcpnt := make(Recipient, l)
	for i := 0; i < l; i++ {
		key := colNames[i]
		value := strings.TrimSpace(fields[i])
		if re, found := slctrs[key]; found {
			if !re.MatchString(value) {
				return nil, ErrNoMatch
			}
		}
		if key == EMailColName {
			switch {
			case len(value) == 0:
				return nil, ErrNoMatch
			case !rxEmail.MatchString(value):
				return nil, ErrInvalidEMail
			}
		}
		rcpnt[key] = value
	}

	return &rcpnt, nil
}

// Recipient maps the column names to their value
type Recipient map[string]string

// Get returns the recipients value for a column name.
func (rcpnt Recipient) Get(colName string) string {
	return rcpnt[colName]
}

// Set sets the value for key
func (rcpnt Recipient) Set(colName, value string) {
	rcpnt[strings.TrimSpace(colName)] = value
}
