package main

import (
	"database/sql"
	"strings"
	// "database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/namsral/flag"
)

type Measurement struct {
	Id   string
	Key  string
	Time sql.NullInt64
	// DateTime time.Time
	Type  string
	Label string
	// Value float64
	Value sql.NullFloat64
}

type argsConfig struct {
	Args     []string
	NtfyBase *string
	TopicId  *string
	DbFile   *string
}

func (ac *argsConfig) valid() bool {
	valid := true
	if *ac.TopicId == "" && *ac.DbFile == "" {
		valid = false
	}
	return valid
}

func (ac *argsConfig) ntfyURL() string {
	return fmt.Sprintf("%s/%s", *ac.NtfyBase, *ac.TopicId)
}

func parseArgs(args []string) *argsConfig {
	ap := &argsConfig{
		Args:     args,
		NtfyBase: flag.String("ntfybase", "https://ntfy.sh", "base domain for Nfty"),
		TopicId:  flag.String("topicid", "", "the topic id from ntfy"),
		DbFile:   flag.String("db", "", "path to database file"),
	}
	flag.String(flag.DefaultConfigFlagname, "", "path to config file")
	flag.Parse()
	// fmt.Printf("db: %+v\n", *ap.DbFile)
	// fmt.Printf("topic: %+v\n", *ap.TopicId)
	if !ap.valid() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	return ap
}

var query = `
with
  tzinfo as (
    select '-04:00' as diff
  ),
  ms_extra as (
    select key,
        strftime('%Y%m%d', datetime(time, 'unixepoch', tzinfo.diff)) as daylocal,
        time, id, label, value, type
    from measurements mm, tzinfo
    where time > strftime('%s',datetime('now',tzinfo.diff,'start of day', '+1 second'))
  ),
  ms as (
    select mm.*,
        min(value) over (partition by daylocal, id) as min_value,
        max(value) over (partition by daylocal, id) as max_value
    from ms_extra as mm
  )
select
    mm.time,
    max(mm.value - mm.min_value) as value
from ms as mm
-- where label = 'electric-recv'
where label = $1
order by mm.time asc
`

func sendNtfy(url, title, message string) {
	// fmt.Printf("%s\n%s\n\n%s\n\n", url, title, message)

	req, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		log.Printf("Error was:\n")
		log.Fatal(err)
	}
	req.Header.Set("Title", title)
	req.Header.Set("Tags", "tada,partying_face")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error posting to ntfy:\n")
		log.Println(err)
	}
	log.Printf("ntfy: %s\n", res.Status)
}

func main() {
	ac := parseArgs(os.Args)

	dbConnStr := fmt.Sprintf("%s?_pragma=journal_mode(WAL)", *ac.DbFile)
	db, err := sqlx.Connect("sqlite", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}

	measurement := Measurement{}
	err = db.Get(&measurement, query, "electric-recv")
	if err != nil {
		log.Printf("SQL Error:\n")
		log.Fatal(err)
	}
	// fmt.Printf(">>> measurement from get: %+v\n", measurement)
	if measurement.Time.Valid && measurement.Value.Float64 >= 3.0 {
		msg := fmt.Sprintf("KWh sent: %.f", measurement.Value.Float64)
		sendNtfy(ac.ntfyURL(), "Sent to Apex ☀️", msg)
	}

}
