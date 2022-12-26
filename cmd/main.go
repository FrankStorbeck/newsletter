package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type config struct {
	from                 string        // email-address of the originator
	infLog               *log.Logger   // stream for logging info
	maxRcpnts            int           // maximum number of newsletters to be sent
	pathToAuthFile       string        // path to auth file
	pathToRecipientsFile string        // path to csv file with data of recipients
	pathToSelectorsFile  string        // path to file with selection criteria for recipients
	skipRcpnts           int           // skip number of recipients before sending newsletters
	sleepTime            time.Duration // time to sleep between sending 2 successive newsletters
	subject              string        // subject of the mail with the newsletter
}

func main() {
	cfg := config{}

	flag.StringVar(&cfg.pathToAuthFile, "auth", ".auth.txt",
		"Path to auth file")
	flag.StringVar(&cfg.from, "from", "",
		"Reply address")
	flag.IntVar(&cfg.maxRcpnts, "max", 100,
		"Maximum number of newsletters to be sent")
	quota := flag.Int("quota", 0,
		"Maximum number of newsletters to be sent during one hour")
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
	usage := flag.Bool("usage", false,
		"Show usage and exit")
	version := flag.Bool("version", false,
		"Show version number and exit")
	flag.Parse()

	if *version {
		fmt.Printf("version 3.0\n")
		return
	}

	if *usage {
		fmt.Printf(man, filepath.Base(os.Args[0])) // for man: see vars.go
		return
	}

	cfg.infLog = log.New(os.Stdout, "      ", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime)

	if len(flag.Args()) <= 0 {
		errLog.Printf("template missing")
		os.Exit(1)
	}

	if cfg.maxRcpnts <= 0 {
		errLog.Fatalf("non positive number for max: %d", cfg.maxRcpnts)
	}

	cfg.sleepTime = 0
	if *quota > 0 {
		cfg.sleepTime = time.Duration((3600 / *quota) * int(time.Second))
	}

	status := 0
	if err := cfg.sendNewsletters(flag.Args()[0]); err != nil {
		errLog.Printf("%s", err.Error())
		status = 1
	}
	os.Exit(status)
}
