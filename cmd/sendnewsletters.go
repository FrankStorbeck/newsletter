package main

import (
	"bufio"
	"fmt"
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

	sls, err := selectors.Read(cfg.pathToSelectorsFile)
	if err != nil {
		return fmt.Errorf("reading selectors file fails: %w", err)
	}

	rF, err := os.Open(cfg.pathToRecipientsFile)
	if err != nil {
		return fmt.Errorf("reading selectors file fails: %w", err)
	}
	defer rF.Close()

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
			name := fmt.Sprintf("%s %s %s", rcp.Get("FirstName"),
				rcp.Get("MiddleNames"), rcp.Get("FamilyName"))
			email := rcp.Get(selectors.EMail)

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
				To:        []sendmail.NamedAddress{{EMail: email, Name: name}},
				Subject:   cfg.subject,
				PlainText: plainBody,
				HTMLText:  htmlBody,
				// Embedd:     []string{filepath.Join(..., "logo.jpg")}
				// Attachments: ...,
			}
			for n := 0; n < 3; n++ {
				if err = dlr.DialAndSend(sendCfg.BuildMessage()); err == nil {
					break
				}
				cfg.infLog.Printf("(%5d) retry (%d) to send email to %q",
					count, n+1, email)
			}
			if err != nil {
				return fmt.Errorf("(%5d) failed to send email to %q: %w",
					count, email, err)
			}
			cfg.infLog.Printf("(%5d) email sent to %q", count, email)
			if cfg.maxRcpnts--; cfg.maxRcpnts <= 0 {
				break
			}
			if cfg.sleepTime > 0 {
				time.Sleep(cfg.sleepTime)
			}
		}
	}
	return nil
}
