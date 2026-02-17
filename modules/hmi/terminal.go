//modules/hmi/terminal.go

package hmi

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/auth"
)

type Terminal struct {
	Kernel  *core.Kernel
	Auth    *auth.AuthManager
	Session *auth.Session
	Scanner *bufio.Scanner
}

func NewTerminal(k *core.Kernel) *Terminal {
	return &Terminal{
		Kernel:  k,
		Auth:    auth.NewAuthManager(k.Vault),
		Scanner: bufio.NewScanner(os.Stdin),
	}
}

// StartLoginFlow handles the initial barrier
func (t *Terminal) StartLoginFlow() error {
	fmt.Println("\n=== SYSTEM LOCKED ===")
	for {
		fmt.Print("1. Login\n2. Sign Up\nSelect: ")
		t.Scanner.Scan()
		choice := strings.TrimSpace(t.Scanner.Text())

		if choice == "2" {
			t.handleSignup()
			continue
		} else if choice == "1" {
			if err := t.handleLogin(); err == nil {
				return nil // Success
			}
			fmt.Println("Error: Invalid credentials. Try again.")
		}
	}
}

func (t *Terminal) handleLogin() error {
	fmt.Print("Username: ")
	t.Scanner.Scan()
	user := strings.TrimSpace(t.Scanner.Text())

	fmt.Print("Password: ")
	t.Scanner.Scan()
	pass := strings.TrimSpace(t.Scanner.Text())

	session, err := t.Auth.Login(user, pass)
	if err != nil {
		return err
	}
	t.Session = session
	fmt.Printf("\n[ACCESS GRANTED] Welcome, %s (%s).\n", user, session.User.Role)
	return nil
}

func (t *Terminal) handleSignup() {
	fmt.Print("\n[NEW USER REGISTRATION]\nUsername: ")
	t.Scanner.Scan()
	user := strings.TrimSpace(t.Scanner.Text())

	fmt.Print("Password: ")
	t.Scanner.Scan()
	pass := strings.TrimSpace(t.Scanner.Text())

	// Default first user to OWNER, others to OPERATOR
	role := "OPERATOR"
	if isEmpty, _ := t.Auth.Vault.IsEmpty("users"); isEmpty {
		role = "OWNER"
	}

	if err := t.Auth.Signup(user, pass, role); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("User created successfully. Please login.")
	}
}

// RunCommandLoop is the main UI loop
func (t *Terminal) RunCommandLoop() {
	fmt.Println("\nType 'help' for commands or 'exit' to shutdown.")
	
	for {
		// The Prompt
		fmt.Printf("\n%s@AIOS-NODE > ", t.Session.User.Username)
		
		if !t.Scanner.Scan() {
			break
		}
		input := strings.TrimSpace(t.Scanner.Text())
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]
		args := parts[1:]

		switch cmd {
		case "help":
			t.showHelp()
		case "status":
			t.showStatus()
		case "dream":
			t.triggerDream(args)
		case "mode":
			t.changeMode(args)
		case "exit":
			fmt.Println("Logging out...")
			return
		default:
			fmt.Printf("Unknown command: %s\n", cmd)
		}
	}
}

func (t *Terminal) showStatus() {
	ctx := t.Kernel.Runtime
	fmt.Printf("--- SYSTEM VITALS ---\n")
	fmt.Printf("Platform:   %s\n", ctx.Platform.Name)
	fmt.Printf("Trust Lvl:  %.2f\n", t.Kernel.TrustLevel())
	fmt.Printf("Boot Mode:  %s\n", ctx.Boot.Type)
	fmt.Printf("Modules:    %d Active\n", len(t.Kernel.Loader.ActiveModules()))
}

func (t *Terminal) changeMode(args []string) {
	if t.Session.User.Role != "OWNER" {
		fmt.Println("PERMISSION DENIED: Only OWNER can change operation modes.")
		return
	}
	if len(args) < 1 {
		fmt.Println("Usage: mode [AUTONOMOUS|MANUAL|SAFE]")
		return
	}
	fmt.Printf("Requesting transition to %s...\n", args[0])
	// In a real implementation, this would call Kernel.SetMode()
}