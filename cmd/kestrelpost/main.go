package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"
	"kestrelpost/internal/ui"
)

func main() {
	hostKeyPath := os.Getenv("KESTREL_HOST_KEY")
	if hostKeyPath == "" {
		hostKeyPath = ".ssh/kestrel_ed25519"
	}

	s, err := wish.NewServer(
		wish.WithAddress(":2222"),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("SSH listening on :2222")
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Printf("server stopped: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Println("stopping SSH server")
	shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(shCtx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Printf("shutdown: %v", err)
	}
	fmt.Fprintln(os.Stderr, "kestrelpost stopped")
}

func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	_ = sess
	return ui.NewModel(), nil
}
