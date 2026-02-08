//internal/terminal/console.go

package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/api/commands"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

type DevConsole struct {
	kernel *core.Kernel
}

func New(k *core.Kernel) *DevConsole {
	return &DevConsole{kernel: k}
}

// Start opens the interactive shell
func (c *DevConsole) Start() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n--- STRATACORE AIOS DEVELOPER CONSOLE ---")
	fmt.Println("Commands: status, nav [x,y], fault [type], halt, exit")

	for {
		fmt.Print("aios-node> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Split(input, " ")

		switch parts[0] {
		case "status":
			t := c.kernel.Trust
			fmt.Printf("[SYSTEM] Mode: %s | Trust: %.2f | Platform: %s\n", 
				t.OperationMode, t.CurrentScore, c.kernel.EnvConfig.Platform.Final)

		case "nav":
			if len(parts) < 2 { continue }
			c.kernel.ProcessCommand(commands.Task{
				Type: commands.CmdNavigate,
				Params: map[string]interface{}{"destination": parts[1]},
			})

		case "fault":
			// Direct access to simulation engine for stress testing
			logging.Warn("[CONSOLE] Manual Fault Injection Triggered: %s", parts[1])
			// Here you would call c.kernel.SimEngine.InjectFault()

		case "halt":
			c.kernel.ProcessCommand(commands.Task{Type: commands.CmdHalt})

		case "exit":
			os.Exit(0)
		}
	}
}