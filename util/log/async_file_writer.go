package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DayTicker struct {
	stop chan struct{}
	C    <-chan time.Time
}

func NewDayTicker() *DayTicker {
	ht := &DayTicker{
		stop: make(chan struct{}),
	}
	ht.C = ht.Ticker()
	return ht
}

func (ht *DayTicker) Stop() {
	ht.stop <- struct{}{}
}

func (ht *DayTicker) Ticker() <-chan time.Time {
	ch := make(chan time.Time)
	go func() {
		day := time.Now().Day()
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for t := range ticker.C {
			if t.Day() != day {
				ch <- t
				day = t.Day()
			}
		}
	}()
	return ch
}

type AsyncFileWriter struct {
	filePath string
	fd       *os.File
	logger   *log.Logger

	mu        sync.Mutex
	dayTicker *DayTicker
}

func NewAsyncFileWriter(filePath string) *AsyncFileWriter {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		panic(fmt.Sprintf("get file path of logger error. filePath=%s, err=%s", filePath, err))
	}

	return &AsyncFileWriter{
		filePath:  absFilePath,
		dayTicker: NewDayTicker(),
	}
}

func (w *AsyncFileWriter) setLogger() {
	if w.logger == nil {
		flag := log.Ldate | log.Lmicroseconds | log.Lshortfile
		w.logger = log.New(w.fd, "", flag)
	} else {
		w.logger.SetOutput(w.fd)
	}
}

func (w *AsyncFileWriter) initLogFile() error {
	var (
		fd  *os.File
		err error
	)

	realFilePath := w.timeFilePath()
	fd, err = os.OpenFile(realFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	w.fd = fd
	if _, err = os.Lstat(w.filePath); err == nil || os.IsExist(err) {
		err = os.Remove(w.filePath)
		if err != nil {
			return err
		}
	}

	if err = os.Symlink(realFilePath, w.filePath); err != nil {
		return err
	}

	w.setLogger()

	return nil
}

func (w *AsyncFileWriter) rotateFile() {
	select {
	case <-w.dayTicker.C:
		w.mu.Lock()
		if err := w.fd.Close(); err != nil {
			panic(err)
		}
		if err := w.initLogFile(); err != nil {
			panic(err)
		}
		w.mu.Unlock()
	default:
	}
}

func (w *AsyncFileWriter) timeFilePath() string {
	return w.filePath + "." + time.Now().Format("2006-01-02")
}
