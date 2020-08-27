package main

import (
	"strconv"
	"sync"
)

func toInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return -1
	}
	return i
}

func toBool(str string) bool {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return b
}

func toFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return -1
	}
	return f
}

// WaitInfinite func
func WaitInfinite() {
	mutex := sync.Mutex{}
	mutex.Lock()
	mutex.Lock()
}

// func jsonizeFromMap(inputmap map[string]string) string {
// 	if bytebuf, err := json.Marshal(inputmap); err != nil {
// 		panic(err.Error())
// 	} else {
// 		return string(bytebuf)
// 	}
// }
