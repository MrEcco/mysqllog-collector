package main

import (
	"io"
	"sync"

	"github.com/hpcloud/tail"
)

// FileStream struct
type FileStream struct {
	// Public
	Filename string
	Config   *tail.Config
	Callback func() []byte
	// Private
	Initialized bool
}

// Init func
func (s *FileStream) Init() {
	if s.Filename == "" {
		ScrapeDaemonPanic("You must specify filestream filename")
	}
	if s.Initialized {
		return
	}
	if s.Config == nil {
		s.Config = &tail.Config{
			Follow: true,
			Location: &tail.SeekInfo{
				Offset: 0,
				// Whence: io.SeekStart,
				Whence: io.SeekEnd,
			},
		}
	}

	// Callback gap
	s.Callback = func() []byte {
		return nil
	}

	mutex := sync.Mutex{}

	mutex.Lock()

	go func() {
		s.Consumer(&mutex)
	}()

	mutex.Lock()
	mutex.Unlock()
	s.Initialized = true
}

// Consumer func
func (s *FileStream) Consumer(mutex *sync.Mutex) {
	channel := make(chan []byte, 16)

	slowquery, errTailFile := tail.TailFile(s.Filename, *s.Config)
	if errTailFile != nil {
		ScrapeDaemonPanic(errTailFile.Error())
	}

	s.Callback = func() []byte {
		return <-channel
	}

	mutex.Unlock()

	for tailLine := range slowquery.Lines {
		line := tailLine.Text
		if errTailFile != nil {
			if errTailFile == io.EOF {
				break
			} else {
				ScrapeDaemonPanic(errTailFile.Error())
			}
		}
		if errTailFile != nil {
			ScrapeDaemonPanic(errTailFile.Error())
		}

		channel <- []byte(line)
	}
}
