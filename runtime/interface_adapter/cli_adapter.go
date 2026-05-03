//runtime/interface_adapter/cli_adapter.go

package interface_adapter

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	auth "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	"go.uber.org/zap"
)

type CLIAuth struct{}

func NewCLIAuth() auth.AuthInterface {
	return &CLIAuth{}
}
func (c *CLIAuth) Authenticate() error {
	return nil
}

type CLIAdapter struct{}

func (c *CLIAdapter) Start(session *user_setting.UserSession) error {
	fmt.Println("CLI session started:", user_setting.UserIdentity)
	return nil
}

func (c *CLIAdapter) Stop(ctx context.Context) error {
	fmt.Println("CLI session stopped")
	return nil
}

func (c *CLIAuth) StartAuthFlow(am *auth.AuthManager) (*user_setting.UserSession, error) {
	return am.LoginOrSignUpInteractive()
}

func (c *CLIAdapter) Notify(msg string) {
	fmt.Println("[CLI]", msg)
}

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

// ============================================================
// CLI
// ============================================================

func startCLI(ctx context.Context, log *zap.Logger) {
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Print("> ")
				in, err := reader.ReadString('\n')
				if err != nil {
					continue
				}
				log.Info("INPUT", zap.String("cmd", strings.TrimSpace(in)))
			}
		}
	}()
}
