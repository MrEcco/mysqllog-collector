package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// SlowQueryRecord struct just representation of log entity in slowquery log stream
type SlowQueryRecord struct {
	User         string  `json:"user,omitempty"`
	Host         string  `json:"host,omitempty"`
	HostIP       string  `json:"host_ip,omitempty"`
	ThreadID     int     `json:"thread_id,omitempty"`
	Schema       string  `json:"schema,omitempty"`
	CacheHit     bool    `json:"cache_hit,omitempty"`
	QueryTime    float64 `json:"query_time,omitempty"`
	LockTime     float64 `json:"lock_time,omitempty"`
	RowsSent     int     `json:"rows_sent,omitempty"`
	RowsExamined int     `json:"rows_examined,omitempty"`
	RowsAffected int     `json:"rows_affected,omitempty"`
	BytesSent    int     `json:"bytes_sent,omitempty"`
	Query        string  `json:"query,omitempty"`
}

var slowQueryRecord1 *regexp.Regexp // Entry
var slowQueryRecord2 *regexp.Regexp
var slowQueryRecord3 *regexp.Regexp
var slowQueryRecord4 *regexp.Regexp

// SlowQueryStreamInit func
func SlowQueryStreamInit() {
	if slowQueryRecord1 == nil {
		slowQueryRecord1, _ = regexp.Compile(`^# User@Host: (\w+)\[(\w+)\] @ *\[([\d\.\:]*)\]`)
	}

	if slowQueryRecord2 == nil {
		slowQueryRecord2, _ = regexp.Compile(`^# Thread_id: (\d+) +Schema: (\w+) +QC_hit: (\w+)`)
	}

	if slowQueryRecord3 == nil {
		slowQueryRecord3, _ = regexp.Compile(`^# Query_time: ([\d\.]+) +Lock_time: ([\d\.]+) +Rows_sent: (\d+) +Rows_examined: (\d+)`)
	}

	if slowQueryRecord4 == nil {
		slowQueryRecord4, _ = regexp.Compile(`^# Rows_affected: (\d+) +Bytes_sent: (\d+)`)
	}
}

// SlowQueryStreamWorker func compute messages. Callback just return next line in log
func SlowQueryStreamWorker(callback func() []byte) SlowQueryRecord {
	var line []byte

	// Seek record entry
	for {
		line = callback()
		if slowQueryRecord1.Match(line) {
			log.Printf("SKIPPING: %s", string(line))
			break
		}
	}

	ret := SlowQueryRecord{}

	// Parse entry
	for {
		line1Submatch := slowQueryRecord1.FindSubmatch(line)
		ret.User = string(line1Submatch[1])
		ret.Host = string(line1Submatch[2])
		ret.HostIP = string(line1Submatch[3])

		failure := false

		line = callback()
		if !slowQueryRecord2.Match(line) {
			failure = true
		} else {
			line2Submatch := slowQueryRecord2.FindSubmatch(line)
			ret.ThreadID = toInt(string(line2Submatch[1]))
			ret.Schema = string(line2Submatch[2])
			ret.CacheHit = toBool(string(line2Submatch[3]))
		}

		line = callback()
		if !slowQueryRecord3.Match(line) {
			failure = true
		} else {
			line3Submatch := slowQueryRecord3.FindSubmatch(line)
			ret.QueryTime = toFloat(string(line3Submatch[1]))
			ret.LockTime = toFloat(string(line3Submatch[2]))
			ret.RowsSent = toInt(string(line3Submatch[3]))
			ret.RowsExamined = toInt(string(line3Submatch[4]))
		}

		line = callback()
		if !slowQueryRecord4.Match(line) {
			fmt.Println("slowQueryRecord4 not match")
			failure = true
		} else {
			line4Submatch := slowQueryRecord4.FindSubmatch(line)
			ret.RowsAffected = toInt(string(line4Submatch[1]))
			ret.BytesSent = toInt(string(line4Submatch[2]))
		}

		query := make([]string, 0)
		for {
			line = callback()
			if slowQueryRecord1.Match(line) {
				break
			}
			query = append(query, string(line))
		}

		if !failure {
			ret.Query = strings.Join(query, "\n")
			RecordPush(ret)
		}
	}
}
