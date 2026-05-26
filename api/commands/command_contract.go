//api/commands/command_contract.go

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

type IncomingCommand struct {
	ID        string                 `json:"id"`
	Type      CommandType            `json:"type"`
	Params    map[string]interface{} `json:"params"`
	Priority  int                    `json:"priority"` // 0 (Critical) to 10 (Low)
	CreatedAt time.Time              `json:"created_at"`
}
