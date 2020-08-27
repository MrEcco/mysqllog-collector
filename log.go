package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// ScrapeDaemonMessage struct
type ScrapeDaemonMessage struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

// Init func fill some fields
func (m *ScrapeDaemonMessage) Init() {
	m.Type = "scrapedaemon"
}

// ScrapeDaemonError func just correctly put internal log message to publish stream
func ScrapeDaemonError(message string) {
	err := ScrapeDaemonMessage{
		Message: message,
	}
	err.Init()

	if pushChannelInitialized {
		RecordPush(err)
	}
}

// ScrapeDaemonPanic func just correctly put internal log message to publish stream and panic
func ScrapeDaemonPanic(message string) {
	ScrapeDaemonError(message)
	os.Exit(1)
}

var pushChannelInitialized bool
var recordPushChannel chan interface{}

// RecordPushInit func prepare record publisher stream
func RecordPushInit() error {
	recordPushChannel = make(chan interface{}, 16)

	go RecordPushConsume(recordPushChannel)

	pushChannelInitialized = true

	return nil
}

// RecordPush func just put record to stdout of container
func RecordPush(record interface{}) {
	recordPushChannel <- record
}

// RecordPushConsume func is main func for pushed records
func RecordPushConsume(queue <-chan interface{}) {
	for {
		Record, errMarshal := json.Marshal(<-queue)
		if errMarshal != nil {
			ScrapeDaemonPanic(errMarshal.Error())
		}
		fmt.Println(string(Record))
	}
}
