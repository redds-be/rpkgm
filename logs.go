package main

import (
	"errors"
	"github.com/pborman/getopt/v2"
	"io"
	"log"
	"os"
)

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, errors.New("could not open log file")
	}
	return logFile, nil
}

func confLog(logFile *os.File) {
	if getopt.GetValue("verbose") == "true" {
		stdoutAndLogFile := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(stdoutAndLogFile)
	} else {
		log.SetOutput(logFile)
	}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func closeLogs(logFile *os.File) {
	err := logFile.Close()
	if err != nil {
		log.Printf("Could not close the log file: '%s': %s", logFile.Name(), err)
		os.Exit(1)
	}
}

func handleLogs(path string) *os.File {
	logFile, err := openLogFile(path)
	if err != nil {
		log.Printf("Could not open the log file '%s': %s", path, err)
		os.Exit(1)
	}
	confLog(logFile)
	return logFile
}
