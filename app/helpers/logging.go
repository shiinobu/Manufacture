package helpers

import (
	"log"
	"os"
)

type Loggers struct {
	INFO  *log.Logger
	ERROR *log.Logger
}

func Logs() (*Loggers, error) {
	x := &Loggers{}
    logDir := "logs"
	logFilePath := logDir + "/error.log"

    // Ensure logs directory exists
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
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
        return nil, err
    }

    // Set up loggers
    x.INFO = log.New(logFile, "INFO: ", log.Ldate|log.Ltime)
    x.ERROR = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime)
    return x, nil
}