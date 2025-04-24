package helpers

import (
	"log"
	"os"
)

var (
	INFO  *log.Logger
	ERROR *log.Logger
)

func Logs() error {
	logDir := "logs"
	logFilePath := logDir + "/error.log"

	// Ensure logs directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil
	}

	// Check if log file exists
	_, err := os.Stat(logFilePath)
	fileExists := !os.IsNotExist(err)

	// Set open flags: append and write-only, add create only if file doesn't exist
	flags := os.O_APPEND | os.O_WRONLY
	if !fileExists {
		flags |= os.O_CREATE
	}

	// Open log file
	logFile, err := os.OpenFile(logFilePath, flags, 0666)
	if err != nil {
		return nil
	}

	// Set up loggers
	return Setup(logFile)
}

func Setup(logfile *os.File) error {
	INFO = log.New(logfile, "INFO: ", log.Ldate|log.Ltime)
	ERROR = log.New(logfile, "ERROR: ", log.Ldate|log.Ltime)
	return nil
}
