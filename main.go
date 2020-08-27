package main

import (
	"flag"
)

func main() {

	// Flags
	// generalFilename := flag.String("general", "", "Path to file with general mysql log (specified by \"general_log_file\")")
	slowqueryFilename := flag.String("slowquery", "", "Path to file with general slowquery log (specified by \"slow_query_log_file\")")
	flag.Parse()

	// Initialize record push stream
	RecordPushInit()

	// if *generalFilename == "" && *slowqueryFilename == "" {
	if *slowqueryFilename == "" {
		ScrapeDaemonPanic("No one log stream specified. Exiting")
	}

	if *slowqueryFilename != "" {
		// Init file stream
		filestream := FileStream{
			Filename: *slowqueryFilename,
		}
		filestream.Init()

		// Init slowquery stream
		SlowQueryStreamInit()

		// Bind streams between
		go SlowQueryStreamWorker(filestream.Callback)
	}

	// if *generalFilename != "" {
	// 	// Init file stream
	// 	filestream := FileStream{
	// 		Filename: *generalFilename,
	// 	}
	// 	filestream.Init()

	// 	// Init general stream
	// 	GeneralStreamInit()

	// 	// Bind streams between
	// 	go GeneralStreamWorker(filestream.Callback)
	// }

	// Just wait
	WaitInfinite()
}
