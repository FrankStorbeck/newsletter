package main

import (
	"encoding/csv"
	"fmt"
	"io"
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

	recNo := 0
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
		recNo++

		if firstLine {
			colNames = make([]string, len(record))
			for i, name := range record {
				colNames[i] = strings.ToLower(strings.TrimSpace(name))
			}
			firstLine = false
			continue
		}

		if recNo <= cfg.skipRcpnts {
			if recNo == cfg.skipRcpnts {
				cfg.infLog.Printf("skipped %d subscribers", recNo)
			}
			continue
		}

		rcp, err := sls.Select(record, colNames)
		if err != nil {
			if err != selectors.ErrNoMatch {
				cfg.errLog.Printf("record %d: %s", recNo, err.Error())
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
			if err == nil {
				break
			}
			cfg.errLog.Printf("record %d: %sretry (%d) to send email to %q",
				recNo, dryTxt, n+1, email)
		}
		if err != nil {
			cfg.errLog.Printf("(%5d) %sfailed to send email to %q: %s",
				recNo, dryTxt, email, err)
		} else {
			cfg.infLog.Printf("record %d: %semail sent to %q", recNo, dryTxt, email)
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
	// // in a dry run generate a error in about 1 of 3 cases
	// if rand.Intn(3) >= 2 {
	// 	err = errors.New("dry run error")
	// }
	return
}
