package logger

import (
	"log"
	"os"
)

var (
	Log *log.Logger
)

func init() {
	file, err := os.OpenFile("discovery.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	Log = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}
