package log

import "fmt"

func Debug(format string, v ...interface{}) {
	w.rotateFile()
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger.Output(2, fmt.Sprintf("[DEBUG] "+format, v...))
}

func Info(format string, v ...interface{}) {
	w.rotateFile()
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger.Output(2, fmt.Sprintf("[INFO] "+format, v...))
}

func Error(format string, v ...interface{}) {
	w.rotateFile()
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger.Output(2, fmt.Sprintf("[ERROR] "+format, v...))
}
