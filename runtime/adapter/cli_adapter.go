//runtime/adapter/cli_adapter.go

package cli

import (
	"context"
	"fmt"
)

type CLIModule struct{}

func NewCLIModule() *CLIModule {
	return &CLIModule{}
}

func (c *CLIModule) Start(ctx context.Context) error {
	fmt.Println("[CLI] started")

	<-ctx.Done()
	return nil
}

func (c *CLIModule) Stop(ctx context.Context) error {
	fmt.Println("[CLI] stopped")
	return nil
}

func (c *CLIModule) Name() string {
	return "cli"
}

func (c *CLIModule) Init(ctx context.Context) error {
	return nil
}

func (c *CLIModule) Health() error {
	return nil
}
