package sendmail

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// Auth holds the data to authorise access to an SMTP mail server
type Auth struct {
	path   string            // File path to file for storing the data
	values map[string]string // Keys and their values
}

// NewAuth returns a new Auth struct.
func NewAuth(path string) *Auth {
	return &Auth{
		path: path,
	}
}

func (auth *Auth) read() ([]byte, error) {
	f, err := os.Open(auth.path)
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return []byte{}, err
	}

	size := stat.Size()
	b := make([]byte, size)
	n, err := f.Read(b)
	if int64(n) != size {
		err = fmt.Errorf("%s incompletely read: %w", auth.path, err)
	}
	return b, err
}

// ReadAuth returns a new Auth struct filled with data stored in file located at
// `filePath`. Each line should have a key followed by a colon and then the
// value. Leading and trailing space characters will be removed from the key as
// well as from their value. The file should hold values for keys "hostname",
// "port", "username", "password" and "sender".
func ReadAuth(filePath string) (*Auth, error) {
	auth := NewAuth(filePath)

	b, err := auth.read()
	if err != nil {
		return auth, err
	}

	auth.values, err = scan(b)
	return auth, err
}

func scan(b []byte) (map[string]string, error) {
	values := make(map[string]string)
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return values, err
		}

		fields := strings.SplitN(scanner.Text(), ":", 2)
		if len(fields) > 1 {
			if k := strings.TrimSpace(fields[0]); len(k) > 0 {
				if v := strings.TrimSpace(fields[1]); len(v) > 0 {
					values[strings.ToLower(k)] = v
				}
			}
		}
	}
	return values, nil
}

// Value returns the value for the key 'k'. Valid keys are "hostname", "port",
// "username", "password" and "sender".
func (auth *Auth) Value(k string) string {
	k = strings.ToLower(k)
	switch k {
	case "hostname", "port", "username", "password", "sender":
		return auth.values[k]
	}
	return ""
}
