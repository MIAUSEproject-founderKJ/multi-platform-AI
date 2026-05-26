//api/hmi/hmi_contract.go - keep contracts only

package hmi

// Update represents a visual command sent to the HUD
type Update struct {
	Source      string      `json:"source"`
	Component   string      `json:"component"`
	Value       interface{} `json:"value"`
	Restart     bool        `json:"restart,omitempty"`
	Defreeze    bool        `json:"defreeze,omitempty"`
	Sensitivity int         `json:"sensitivity,omitempty"`
	TBC         bool        `json:"tbc,omitempty"`
}

// DisplayProvider allows different UIs (OpenGL, Web, Terminal) to render the AI state
type DisplayProvider interface {
	PushUpdate(upd Update)
	GetStatus() string
	TBC() bool
}
