//internal/logging/logger.go

package logging

import (
	"log"
)

func Debug(format string, v ...interface{}) { log.Printf("[DEBUG] "+format, v...) }
func Info(format string, v ...interface{})  { log.Printf("[INFO] "+format, v...) }
func Warn(format string, v ...interface{})  { log.Printf("[WARN] "+format, v...) }
func Error(format string, v ...interface{}) { log.Printf("[ERROR] "+format, v...) }

// For the Safety Interlock error
func ProtectIncidentLog(err error) {
	log.Printf("[CRITICAL] Incident Log Protected: %v", err)
}
