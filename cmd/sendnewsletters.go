package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"storbeck.nl/newsletter/pkg/selectors"
	"storbeck.nl/newsletter/pkg/sendmail"
)

// sendNewsletters does the work
func (cfg *config) sendNewsletters(pathToTmplt string) error {
	tmplt, err := getTemplate(pathToTmplt, "mail_bodies")
	if err != nil {
		return fmt.Errorf("parsing template fails: %w", err)
	}

	ath, err := sendmail.ReadAuth(cfg.pathToAuthFile)
	if err != nil {
		return fmt.Errorf("reading auth file fails: %w", err)
	}

	dlr, err := sendmail.NewDialer(ath)
	if err != nil {
		return err
	}

	sls, err := selectors.New(cfg.pathToSelectorsFile)
	if err != nil {
		return fmt.Errorf("reading selectors file fails: %w", err)
	}

	f, err := os.Open(cfg.pathToSubscribersFile)
	if err != nil {
		return fmt.Errorf("reading subscribers file fails: %w", err)
	}

	count := 0
	firstLine := true
	var colNames []string

	r := csv.NewReader(f)
	r.Comma = ';'
	r.Comment = '#'
	r.FieldsPerRecord = -1
	r.ReuseRecord = true

	for {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("reading record from subscribers file fails: %w", err)
		}

		if firstLine {
			colNames = make([]string, len(record))
			for i, name := range record {
				colNames[i] = strings.ToLower(strings.TrimSpace(name))
			}
			firstLine = false
			continue
		}

		count++
		if count <= cfg.skipRcpnts {
			if count == cfg.skipRcpnts {
				cfg.infLog.Printf("skipped %d newsletters", count)
			}
			continue
		}

		rcp, err := sls.Select(record, colNames)
		if err != nil {
			if err != selectors.ErrNoMatch {
				cfg.errLog.Printf("(%5d) %s", count, err.Error())
			}
			continue
		}

		email := rcp.Get(selectors.EMailColName) // rcp always contains a valid email here

		plainBody, err := sendmail.Body("plainBody", tmplt, rcp)
		if err != nil {
			return err
		}
		htmlBody, err := sendmail.Body("htmlBody", tmplt, rcp)
		if err != nil {
			return err
		}
		if len(plainBody) == 0 && len(htmlBody) == 0 {
			// tmplt seems to hold only a plain version for the template
			plainBody, err = sendmail.Body("", tmplt, rcp)
			if err != nil {
				return err
			}
		}

		from := ath.Value("sender")
		if len(cfg.from) > 0 {
			from = cfg.from
		}

		sendCfg := sendmail.Config{
			Sender:    ath.Value("sender"),
			From:      sendmail.NamedAddress{EMail: from, Name: ""},
			To:        []sendmail.NamedAddress{{EMail: email, Name: ""}},
			Subject:   cfg.subject,
			PlainText: plainBody,
			HTMLText:  htmlBody,
			// Embedd:     []string{filepath.Join(..., "logo.jpg")}
			// Attachments: ...,
		}
		dryTxt := ""
		if cfg.dry {
			dryTxt = "dry run: "
		}
		for n := 0; n < 3; n++ {
			if !cfg.dry {
				err = dlr.DialAndSend(sendCfg.BuildMessage())
			} else {
				err = dryDialAndSend()
			}
			if err != nil {
				cfg.errLog.Printf("(%5d) %sretry (%d) to send email to %q",
					count, dryTxt, n+1, email)
			}
		}
		if err != nil {
			cfg.errLog.Printf("(%5d) %sfailed to send email to %q: %s",
				count, dryTxt, email, err)
		} else {
			cfg.infLog.Printf("(%5d) %semail sent to %q", count, dryTxt, email)
		}

		if cfg.maxRcpnts--; cfg.maxRcpnts <= 0 {
			break
		}

		if cfg.sleepTime > 0 {
			// obey to quota
			time.Sleep(cfg.sleepTime)
		}
	}

	return nil
}

func dryDialAndSend() (err error) {
	// in a dry run generate a error in about 1 of 20 cases
	if rand.Intn(20) > 18 {
		err = errors.New("dry run error")
	}
	return
}
