//api/commands/translator.go

package commands

import (
	"time"
)

type CommandType string

const (
	CmdNavigate CommandType = "NAVIGATE"
	CmdScan     CommandType = "PERCEPTION_SCAN"
	CmdHalt     CommandType = "EMERGENCY_HALT"
	CmdSync     CommandType = "DATA_SYNC"
)

type Task struct {
	ID        string                 `json:"id"`
	Type      CommandType            `json:"type"`
	Params    map[string]interface{} `json:"params"`
	Priority  int                    `json:"priority"` // 1 (Low) to 10 (Critical)
	CreatedAt time.Time              `json:"created_at"`
}