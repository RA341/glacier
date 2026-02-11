package app

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/ra341/glacier/pkg/logger"
	"github.com/rs/zerolog"
)

const socketFile = "glacier.sock"

const StartUIMsg = "StartUI"

type SocketManager struct {
	socketPath string
	log        *zerolog.Logger
}

func NewSocketManager(logDir string) *SocketManager {
	fullSocket := filepath.Join(os.TempDir(), socketFile)
	logWriter := NewFileLogger(filepath.Join(logDir, "socket"))
	subLog := logger.CreateLogger("debug", false, logWriter).Str("component", "socket").Logger()

	return &SocketManager{
		socketPath: fullSocket,
		log:        &subLog,
	}
}

// will exit if found
func (sm *SocketManager) exitIfAlreadyRunning() {
	conn, err := net.Dial("unix", sm.socketPath)
	if err != nil {
		// no prev socket found this is the only instance
		sm.log.Info().Err(err).Msg("Failed to connect to socket")
		return
	}

	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			sm.log.Error().Err(err).Msg("Failed to close socket")
		}
	}(conn)

	sm.log.Info().Msg("Another instance is running. Sending instructions...")
	message := strings.Join(os.Args[1:], " ")
	if message == "" {
		// launch desktop
		message = StartUIMsg
	}
	_, err = conn.Write([]byte(message))
	if err != nil {
		sm.log.Warn().Err(err).Msg("could not send message")
	}

	os.Exit(0)
}

func (sm *SocketManager) setupSocketHandler(ctx context.Context, startUIFn func()) error {
	_ = os.Remove(sm.socketPath)

	listener, err := net.Listen("unix", sm.socketPath)
	if err != nil {
		return err
	}
	defer listener.Close()

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	sm.log.Info().Str("addr", listener.Addr().String()).Msg("Socket handler listening on")

	handler := func(conn net.Conn) {
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				sm.log.Printf("Error closing socket connection: %v", err)
			}
		}(conn)

		buf, err := io.ReadAll(conn)
		if err != nil {
			sm.log.Error().Err(err).Msg("Error reading from socket")
			return
		}

		msg := string(buf)
		sm.log.Info().Str("msg", msg).Msg("[Primary Instance] Received command from another instance")

		if msg == StartUIMsg {
			// send message to start desktop ui
			startUIFn()
		}
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("could not accept connection %w", err)
		}
		go handler(conn)
	}
}
