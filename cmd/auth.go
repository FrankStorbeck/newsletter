package main

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// values used as key in configuraton files
const (
	from = "from"
	host = "hostname"
	port = "port"
	pwd  = "password"
	user = "username"

	postmaster = "Postmaster"
)

// Auth holds the data to authorise access to the SMTP mail server
type Auth struct {
	From string // Sender used when connecting to the SMTP mail server
	Host string // SMTP host
	Port int    // Port used by SMTP host
	Pwd  string // Password for SMTP host
	User string // Username to login on SMTP host
}

// readAuth returns a struct holding the data in a configuration file.
// Each line should have a key followed by a colon and then the value.
// Leading and trailing space characters will be removed from the key
// and the value.
// Recognised keys are shown in the following example:
//
// hostname: a.mailhost.dom
// port:     365
// username: someone@mailhost.dom
// password: aap je])?5!
// from:     sender@mailhost.dom
//
func readAuth(filePath string) (*Auth, error) {
	d := make(map[string]string)
	auth := &Auth{}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return auth, err
	}

	lines := strings.Split(string(b), "\n")
	for i := 0; i < len(lines); i++ {
		fields := strings.SplitN(lines[i], ":", 2)
		if len(fields) > 1 {
			fields[0] = strings.TrimSpace(fields[0])
			fields[1] = strings.TrimSpace(fields[1])
			if len(fields[0]) > 0 && len(fields[1]) > 0 {
				d[fields[0]] = fields[1]
			}
		}
	}
	auth.From = d[from]
	auth.Host = d[host]
	port, err := strconv.Atoi(d[port])
	if err == nil {
		auth.Port = port
	}
	auth.User = d[user]
	auth.Pwd = d[pwd]
	return auth, nil
}
