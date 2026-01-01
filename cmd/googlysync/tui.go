package main

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/sandeepkv93/googlysync/internal/config"
	"github.com/sandeepkv93/googlysync/internal/ipc"
	ipcgen "github.com/sandeepkv93/googlysync/internal/ipc/gen"
)

type statusMsg struct {
	state   string
	message string
	at      time.Time
}

type errMsg struct {
	err error
}

type model struct {
	socketPath string
	interval   time.Duration
	status     statusMsg
	err        error
	quitting   bool
}

func newModel(socketPath string, interval time.Duration) model {
	return model{
		socketPath: socketPath,
		interval:   interval,
	}
}

func (m model) Init() tea.Cmd {
	return pollStatusCmd(m.socketPath, m.interval)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.status = msg
		m.err = nil
		return m, pollStatusCmd(m.socketPath, m.interval)
	case errMsg:
		m.err = msg.err
		return m, tea.Tick(m.interval, func(time.Time) tea.Msg {
			return pollNowMsg{}
		})
	case pollNowMsg:
		return m, pollStatusCmd(m.socketPath, m.interval)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "r":
			return m, pollStatusCmd(m.socketPath, 0)
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "\n"
	}
	if m.err != nil {
		return fmt.Sprintf("googlysync status\n\nerror: %v\n\nq to quit, r to retry\n", m.err)
	}
	if m.status.at.IsZero() {
		return "googlysync status\n\nloading...\n\nq to quit\n"
	}
	return fmt.Sprintf("googlysync status\n\n%s: %s\nupdated: %s\n\nq to quit, r to refresh\n", m.status.state, m.status.message, m.status.at.Format(time.RFC3339))
}

type pollNowMsg struct{}

func pollStatusCmd(socketPath string, interval time.Duration) tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.NewConfigWithOptions(config.Options{SocketPath: socketPath})
		if err != nil {
			return errMsg{err: err}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		conn, err := ipc.Dial(ctx, cfg.SocketPath)
		if err != nil {
			return errMsg{err: err}
		}
		defer conn.Close()

		client := ipcgen.NewSyncStatusClient(conn)
		resp, err := client.GetStatus(ctx, &ipcgen.Empty{})
		if err != nil {
			return errMsg{err: err}
		}
		if resp == nil || resp.Status == nil {
			return errMsg{err: fmt.Errorf("no status returned")}
		}
		return statusMsg{
			state:   resp.Status.State.String(),
			message: resp.Status.Message,
			at:      time.Now(),
		}
	}
}
