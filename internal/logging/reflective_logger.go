//MIAUSEproject-founderKJ/multi-platform-AI/internal/logging/reflective_logger.go

package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/api/hmi"
)

type ReflectiveLogger struct {
	FileHandle *os.File
	HMIPipe    chan hmi.ProgressUpdate
}

type LogEntry struct {
	Timestamp string `json:"t"`
	Level     string `json:"l"`
	Message   string `json:"m"`
}

// Info logs a message and reflects it to the HUD
func (rl *ReflectiveLogger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "INFO",
		Message:   msg,
	}

	// 1. Write to Disk (JSONL format for AI training)
	jsonData, _ := json.Marshal(entry)
	rl.FileHandle.WriteString(string(jsonData) + "\n")

	// 2. Reflect to HUD
	if rl.HMIPipe != nil {
		rl.HMIPipe <- hmi.ProgressUpdate{
			Stage:   "LOG_REFLECT",
			Message: fmt.Sprintf("[%s] %s", entry.Level, entry.Message),
		}
	}
}
