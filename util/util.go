package util

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/TwiN/go-color"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, color.Green, log.LstdFlags)
}

// Map2String can transfer a string-string map into a raw string
func Map2string(params map[string]string) string {
	var query string
	for k, v := range params {
		query += k + "=" + v + "&"
	}
	query = query[:len(query)-1]
	return query
}

// Map2String can transfer a string-string map into url value struct
func Map2Params(params map[string]string) url.Values {
	value := url.Values{}
	for key, param := range params {
		value[key] = []string{param}
	}
	return value
}

// GetTimestamp can obtain current ts
func GetTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

// StringContains judge whether val exist in array
func IntContain(array []int, val int) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

/*****************************************
 *              Log helpers              *
 *****************************************/
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
