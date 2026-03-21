//modules\cognition_module.go

package modules

import (
	"context"
	"encoding/json"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type CognitionModule struct {
	ctx *schema.RuntimeContext
}

func (m *CognitionModule) Init(ctx *schema.RuntimeContext) error {

	m.ctx = ctx

	ctx.Bus.Subscribe("audio.intent", m.Handle)

	return nil
}

func (m *CognitionModule) Handle(ctx context.Context, payload []byte) error {

	intent := parseIntent(payload)

	if intent.Confidence < 0.75 {
		return nil
	}

	task := plan(intent)

	data, _ := json.Marshal(task)

	return m.ctx.Data.Bus.Publish(ctx, "action.request", data)
}
