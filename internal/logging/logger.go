//internal/logging/logger.go

package logging

import "fmt"

func Info(format string, v ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", v...)
}

func Error(format string, v ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", v...)
}