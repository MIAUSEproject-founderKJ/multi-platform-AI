//cmd/aios/runtime/session.go

/*The session handles:
• External IO
• Lifecycle binding
• Controlled shutdown
• Backpressure*/

package runtime

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/agent"
)

type Session interface {
	Start(context.Context) error
	Stop()
}

type session struct {
	ctx        *ExecutionContext
	agent      *agent.AgentRuntime
	listener   net.Listener
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

func NewSession(ctx *ExecutionContext, agent *agent.AgentRuntime) Session {
	return &session{
		ctx:   ctx,
		agent: agent,
	}
}

func (s *session) Start(parent context.Context) error {

	ctx, cancel := context.WithCancel(parent)
	s.cancelFunc = cancel

	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		return err
	}

	s.listener = ln

	s.wg.Add(1)
	go s.acceptLoop(ctx)

	return nil
}

func (s *session) acceptLoop(ctx context.Context) {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			continue
		}

		s.wg.Add(1)
		go s.handleConnection(ctx, conn)
	}
}

func (s *session) handleConnection(ctx context.Context, conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	buffer := make([]byte, 8192)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return
			}
			continue
		}

		payload := buffer[:n]

		_ = s.agent.Process(ctx, s.ctx.Optimizer, payload)
	}
}

func (s *session) Stop() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	if s.listener != nil {
		s.listener.Close()
	}

	s.wg.Wait()
}

/*Key Properties:

• Each connection handled in its own goroutine
• Context cancellation respected
• Read deadline prevents stuck connections
• Controlled shutdown with WaitGroup*/