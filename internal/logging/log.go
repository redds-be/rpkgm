package logging

import (
	"fmt"
	"log"
	"os"
)

// LogToFile logs given content to a log file.
func LogToFile(format string, toLog ...any) {
	// File mode to use
	const fileMode = 0o666

	// Open the log file or create it if it does not exist
	logFile, err := os.OpenFile("var/log/rpkgm.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, fileMode)
	if err != nil {
		log.Printf("rpkgm could not open the log file. Error: %s\n", err)

		return
	}

	// Set the output as the log file
	log.SetOutput(logFile)

	// Log the given content
	log.Printf(fmt.Sprintf("%v\n", format), toLog...)

	// Close the log file
	err = logFile.Close()
	if err != nil {
		log.Printf("rpkgm could not close the log file. Error: %s\n", err)

		return
	}
}
