package util

import (
	"log"
	"os"

	"github.com/TwiN/go-color"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, color.Green, log.LstdFlags)
}

func Debug(format string, v ...interface{}) {
	format = "[DEBUG]" + format
	log.Default().Printf(format, v...)
}

func Info(format string, v ...interface{}) {
	format = "[INFO]" + format
	log.Default().Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	format = "[ERROR]" + format
	log.Default().Printf(format, v...)
}

func PrintColor(format string, v ...interface{}) {
	logger.Printf(format, v...)
}
