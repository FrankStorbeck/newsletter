package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/gomail.v2"
	"storbeck.nl/newsletter/pkg/selectors"
)

type config struct {
	infLog               *log.Logger // stream for logging info
	maxRcpnts            int         // maximum number of newsletters to be sent
	pathToAuthFile       string      // path to auth file
	pathToRecipientsFile string      // path to csv file with data of recipients
	pathToSelectorsFile  string      // path to file with selection creteria for recipients
	skipRcpnts           int         // skip number of recipients before sending newsletters
	subject              string      // subject of the mail with the newsletter
}

func main() {
	cfg := config{}

	flag.StringVar(&cfg.pathToAuthFile, "auth", ".auth.txt",
		"Path to auth file")
	flag.IntVar(&cfg.maxRcpnts, "max", 100,
		"Maximum number of newsletters to be sent")
	flag.IntVar(&cfg.skipRcpnts, "skip", 0,
		"Number selected recipients to be skipped before sending the 1st newsletter")
	flag.StringVar(&cfg.pathToRecipientsFile, "recipients",
		filepath.Join("model", "recipients.csv"),
		"Path to recipients file")
	flag.StringVar(&cfg.pathToSelectorsFile, "selectors",
		filepath.Join("model", "selectors.txt"),
		"Path to selectors file")
	flag.StringVar(&cfg.subject, "subject", "Newsletter",
		"Subject of the mailing")
	test := flag.Bool("test", false,
		"Test by sending the newsletter only to some selected adresses")
	usage := flag.Bool("usage", false,
		"Show usage and exit")
	version := flag.Bool("version", false,
		"Show version number and exit")
	flag.Parse()

	if *version {
		fmt.Printf("version 2.0\n")
		return
	}

	if *usage {
		fmt.Printf(man, filepath.Base(os.Args[0])) // for man: see vars.go
		return
	}

	cfg.infLog = log.New(os.Stdout, "      ", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime)

	if cfg.maxRcpnts <= 0 {
		errLog.Fatalf("non positive number for max: %d", cfg.maxRcpnts)
	}

	if *test {
		cfg.pathToSelectorsFile = filepath.Join(".", "tests", "selectors_test.txt")
		cfg.pathToRecipientsFile = filepath.Join(".", "tests", "recipients_test.csv")
	}

	status := 0
	if err := cfg.sendNewsletters(flag.Args()); err != nil {
		errLog.Printf("%s", err.Error())
		status = 1
	}
	os.Exit(status)
}

// sendNewsletters does the work
func (cfg *config) sendNewsletters(tmplts []string) error {
	parsedTmplts, err := parseTemplates(tmplts)
	if err != nil {
		return fmt.Errorf("parsing templates fails: %w", err)
	}

	ext := []string{".txt", ".html"}
	validTemplt := false
	for _, e := range ext {
		if parsedTmplts[e] != nil {
			validTemplt = true
			break
		}
	}
	if !validTemplt {
		return fmt.Errorf("no valid templates found")
	}

	a, err := readAuth(cfg.pathToAuthFile)
	if err != nil {
		return fmt.Errorf("reading auth file fails: %w", err)
	}

	sls, err := selectors.Read(cfg.pathToSelectorsFile)
	if err != nil {
		return fmt.Errorf("selectors.Read() returns an error: %w", err)
	}

	rF, err := os.Open(cfg.pathToRecipientsFile)
	if err != nil {
		return fmt.Errorf("os.Open() returns an error: %w", err)
	}
	defer rF.Close()

	d := gomail.NewDialer(a.Host, a.Port, a.User, a.Pwd)
	sC, err := d.Dial()
	if err != nil {
		return fmt.Errorf("dialing to %s fails: %w", a.Host, err)
	}
	defer sC.Close()

	m := gomail.NewMessage()
	count := 0
	scanner := bufio.NewScanner(rF)

	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return fmt.Errorf("reading line from recipients file fails: %w", err)
		}
		line := scanner.Text()
		if n := strings.Index(line, "#"); n >= 0 {
			line = line[:n]
		}
		if len(line) == 0 {
			continue
		}

		rcp, ok := sls.TestRecipient(scanner.Text())
		if ok {
			count++
			if count <= cfg.skipRcpnts {
				if count == cfg.skipRcpnts {
					cfg.infLog.Printf("skipped %d newsletters (%q)",
						count, "..., "+rcp.Get(selectors.EMail))
				}
				continue
			}

			m.SetHeader("From", a.From)
			name := fmt.Sprintf("%s %s %s", rcp.Get("FirstName"),
				rcp.Get("MiddleNames"), rcp.Get("FamilyName"))
			m.SetAddressHeader("To", rcp.Get(selectors.EMail), name)
			m.SetHeader("Subject", cfg.subject)

			for i, e := range ext {
				if parsedTmplts[e] != nil {
					body := new(bytes.Buffer)
					err = parsedTmplts[e].Execute(body, rcp)
					if err != nil {
						return fmt.Errorf("executing template fails: %w", err)
					}
					s := body.String()
					if i == 0 {
						// e is equal to ".txt"
						m.SetBody("text/plain", s)
					} else {
						// e is equal to ".html"
						m.AddAlternative("text/html", s)
					}
				}
			}

			if err := gomail.Send(sC, m); err != nil {
				return fmt.Errorf("(%5d) failed to send email to %q: %w",
					count, rcp.Get(selectors.EMail), err)
			}
			cfg.infLog.Printf("(%5d) mailing sent to %q",
				count, rcp.Get(selectors.EMail))

			m.Reset()
			if cfg.maxRcpnts--; cfg.maxRcpnts <= 0 {
				break
			}
		}
	}
	return nil
}
