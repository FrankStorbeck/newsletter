package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

type config struct {
	dry                   bool          // run in dry mode
	errLog                *log.Logger   // stream for logging errors
	from                  string        // email-address of the originator
	infLog                *log.Logger   // stream for logging info
	maxRcpnts             int           // maximum number of newsletters to be sent
	pathToAuthFile        string        // path to auth file
	pathToSubscribersFile string        // path to csv file with data of subscribers
	pathToSelectorsFile   string        // path to file with selection criteria for recipients
	skipRcpnts            int           // skip number of recipients before sending newsletters
	sleepTime             time.Duration // time to sleep between sending 2 successive newsletters
	subject               string        // subject of the mail with the newsletter
}

func main() {
	cfg := config{}

	flag.BoolVar(&cfg.dry, "dry", false,
		"Do dry run")
	flag.StringVar(&cfg.pathToAuthFile, "auth", ".auth.txt",
		"Path to auth file")
	flag.StringVar(&cfg.from, "from", "",
		"Reply address")
	flag.IntVar(&cfg.maxRcpnts, "max", 0,
		"Maximum number of newsletters to be sent")
	quota := flag.Int("quota", 0,
		"Maximum number of newsletters to be sent during one hour")
	flag.StringVar(&cfg.pathToSelectorsFile, "selectors", "selectors.txt",
		"Path to selectors file")
	flag.IntVar(&cfg.skipRcpnts, "skip", 0,
		"Number selected recipients to be skipped before sending the 1st newsletter")
	flag.StringVar(&cfg.subject, "subject", "News letter",
		"Subject of the mailing")
	flag.StringVar(&cfg.pathToSubscribersFile, "subscribers", "subscribers.csv",
		"Path to recipients file")
	version := flag.Bool("version", false,
		"Show version number and exit")
	flag.Parse()

	if *version {
		fmt.Printf("sendnewsletter v4.0\n")
		return
	}

	cfg.infLog = log.New(os.Stdout, "      ", log.Ldate|log.Ltime)
	cfg.errLog = log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime)

	if len(flag.Args()) <= 0 {
		cfg.errLog.Printf("template missing")
		os.Exit(1)
	}

	if cfg.maxRcpnts <= 0 {
		cfg.maxRcpnts = math.MaxInt64
	}

	cfg.sleepTime = 0
	if *quota > 0 {
		cfg.sleepTime = time.Duration((3600 / *quota) * int(time.Second))
	}

	status := 0
	if err := cfg.sendNewsletters(flag.Args()[0]); err != nil {
		cfg.errLog.Printf("%s", err.Error())
		status = 1
	}
	os.Exit(status)
}
