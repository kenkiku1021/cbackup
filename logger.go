package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"path/filepath"
)

var logger *log.Logger

func InitLogger(cfg Config) {
	logfileName := filepath.Join(cfg.LogDir, getLogfileName())
	logfile, err := os.Create(logfileName)
	if err != nil {
		log.Fatal("Cannot create log file (%v): %v", logfileName, err)
	}
	logger = log.New(logfile, "cbackup: ",
		log.Ldate | log.Ltime | log.Lshortfile)
}

func getLogfileName() string {
	t := time.Now()
	return fmt.Sprintf("cbackup-%v.log", t.Format("20160102-150405"))
}

