package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"

	"fyne.io/systray"
	"github.com/natefinch/lumberjack"
	"github.com/ra341/glacier/frost/app/icon"
	"github.com/ra341/glacier/frost/config"
	"github.com/ra341/glacier/pkg/logger"
	"github.com/ra341/glacier/shared/api"
	sharedConfig "github.com/ra341/glacier/shared/config"
	"github.com/rs/zerolog"
)

type Tray struct {
	conf       *sharedConfig.Service[config.Config]
	serverOpts []api.ServerOpt

	ctx    context.Context
	cancel context.CancelFunc

	wg            sync.WaitGroup
	uiRunning     atomic.Bool
	serverRunning atomic.Bool

	trayLog *zerolog.Logger
}

func NewDesktop(opts ...api.ServerOpt) {
	// initial config for tray launch will be overwritten by server init
	conf := config.New(false)

	sm := NewSocketManager(conf.Get().Files.LogsDir)
	sm.exitIfAlreadyRunning()

	t := Tray{
		serverOpts: opts,
		conf:       conf,
	}
	l := t.makeLogger("tray")
	t.trayLog = &l

	go func() {
		err := sm.setupSocketHandler(context.Background(), t.StartUI)
		if err != nil {
			t.trayLog.Fatal().Err(err).Msg("Failed to start socket listener")
			return
		}
	}()

	t.Start()
}

func (t *Tray) Start() {
	t.startServices()
	systray.Run(t.onReady, t.onExit)
}

func (t *Tray) startServices() {
	// hold until functions clean up
	t.wg.Wait()

	// then reset context
	ctx, cancel := context.WithCancel(context.Background())
	t.ctx = ctx
	t.cancel = cancel

	t.wg.Go(t.startServer)

	if t.conf.Get().Desktop.StartSilent {
		t.trayLog.Info().Msg("UI set to launch silently")
		return
	}
	t.wg.Go(t.startUI)
}

func (t *Tray) startServer() {
	if t.serverRunning.Load() {
		t.trayLog.Info().Msg("Server is already running")
		return
	}

	t.trayLog.Info().Msg("Starting server...")

	defer func() {
		t.trayLog.Info().Msg("Server stopped")
		t.serverRunning.Store(false)
	}()
	t.serverRunning.Store(true)

	app := New()
	t.conf = app.Conf

	finalOpts := append(t.serverOpts, api.WithCtx(t.ctx))
	StartServerRaw(app, finalOpts...)
}

func (t *Tray) StartUI() {
	go t.startUI()
}

func (t *Tray) startUI() {
	subLog := t.makeLogger("ui")

	if t.uiRunning.Load() {
		subLog.Info().Msg("UI is running")
		return
	}

	subLog.Info().Msg("Starting UI")
	defer func() {
		t.uiRunning.Store(false)
	}()
	t.uiRunning.Store(true)

	exePath := "ui/ui"
	if runtime.GOOS == "windows" {
		exePath += ".exe"
	}

	err := NewUI(t.ctx, exePath, &subLog)
	if err != nil {
		if errors.Is(t.ctx.Err(), context.Canceled) {
			subLog.Debug().Msg("UI closed by user")
			return
		}

		ShowErr(fmt.Sprintf("Failed to start UI: %v", err))
	}
}

func (t *Tray) makeLogger(name string) zerolog.Logger {
	get := t.conf.Get()

	logWriter := NewFileLogger(filepath.Join(get.Files.LogsDir, name))
	return logger.
		CreateLogger(
			get.Logger.Level,
			get.Logger.Verbose,
			logWriter,
		).
		Str("component", name).Logger()
}

func NewFileLogger(filename string) io.Writer {
	logFile := &lumberjack.Logger{
		Filename:   filename + ".log",
		MaxSize:    3,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
	return logFile
}

func (t *Tray) onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Frost")
	systray.SetTooltip("Frost")
	systray.SetOnTapped(func() {
		t.StartUI()
	})

	mUI := systray.AddMenuItem("Open", "Start the UI")
	mLaunchWebUI := systray.AddMenuItem("Open in browser", "Launched the Web UI")

	systray.AddSeparator()

	logsDir := systray.AddMenuItem("Logs", "Open Logs directory")

	systray.AddSeparator()

	mServer := systray.AddMenuItem("Restart", "Restart the app")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		for {
			select {
			case <-logsDir.ClickedCh:
				openFolder(t.conf.Get().Files.LogsDir)
			case <-mLaunchWebUI.ClickedCh:
				//t.trayLog.Info().Msg("launching web browser")
				openURL(fmt.Sprintf(
					"http://localhost:%d",
					t.conf.Get().Server.Port,
				))
			case <-mServer.ClickedCh:
				t.trayLog.Info().Msg("restarting frost")
				t.cancel()
				t.startServices()
			case <-mUI.ClickedCh:
				//t.trayLog.Info().Msg("starting UI")
				go t.startUI()
			case <-mQuit.ClickedCh:
				t.trayLog.Info().Msg("exiting frost")
				t.cancel()
				systray.Quit()
			}
		}
	}()
}

func (t *Tray) onExit() {
	t.cancel()
}

func openFolder(path string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin": // macOS
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		ShowErr(fmt.Sprintf("unsupported platform: %s", runtime.GOOS))
	}

	err := cmd.Start()
	if err != nil {
		ShowErr(fmt.Sprintf("Failed to open folder: %v", err))
	}
}

func openURL(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		ShowErr(fmt.Sprintf("unsupported platform: %s", runtime.GOOS))
	}

	err := cmd.Start()
	if err != nil {
		ShowErr(fmt.Sprintf("could not launch URL: %v", err))
	}
}
