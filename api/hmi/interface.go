//api/hmi/interface.go

package hmi

// Update represents a visual command sent to the HUD
type Update struct {
	Source    string      `json:"source"`
	Component string      `json:"component"`
	Value     interface{} `json:"value"`
}

// DisplayProvider allows different UIs (OpenGL, Web, Terminal) to render the AI state
type DisplayProvider interface {
	PushUpdate(upd Update)
	GetStatus() string
}