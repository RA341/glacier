package app

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"

	"fyne.io/systray"
	"github.com/ra341/glacier/frost/app/icon"
	"github.com/ra341/glacier/frost/config"
	"github.com/ra341/glacier/shared/api"
	"github.com/rs/zerolog/log"
)

type Tray struct {
	ctx    context.Context
	cancel context.CancelFunc

	conf *config.Service

	wg            sync.WaitGroup
	uiRunning     atomic.Bool
	serverRunning atomic.Bool
	serverOpts    []api.ServerOpt
}

func NewDesktop(opts ...api.ServerOpt) {
	sm := NewSocketManager()
	sm.exitIfAlreadyRunning()

	t := Tray{
		serverOpts: opts,
		conf:       config.New(false), // initial config for UI launch will be overwritten by server init
	}

	go func() {
		err := sm.setupSocketHandler(context.Background(), t.StartUI)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start socket listener")
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
		log.Info().Msg("UI set to launch silently")
		return
	}
	t.wg.Go(t.startUI)
}

func (t *Tray) startServer() {
	if t.serverRunning.Load() {
		log.Info().Msg("Server is already running")
		return
	}

	log.Info().Msg("Starting server...")

	defer func() {
		log.Info().Msg("Server stopped")
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
	if t.uiRunning.Load() {
		log.Info().Msg("UI is running")
		return
	}

	log.Info().Msg("Starting UI")
	defer func() {
		t.uiRunning.Store(false)
	}()
	t.uiRunning.Store(true)

	exePath := "ui/ui"
	if runtime.GOOS == "windows" {
		exePath += ".exe"
	}

	err := NewUI(t.ctx, exePath)
	if err != nil {
		if errors.Is(t.ctx.Err(), context.Canceled) {
			log.Info().Msg("Process stopped by user")
			return
		}

		ShowErr(fmt.Sprintf("Failed to start UI: %v", err))
	}
}

func (t *Tray) onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Frost")
	systray.SetTooltip("Frost")

	mUI := systray.AddMenuItem("Open", "Start the UI")
	mLaunchWebUI := systray.AddMenuItem("Open in browser", "Launched the Web UI")
	mServer := systray.AddMenuItem("Restart", "Restart the app")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		for {
			select {
			case <-mLaunchWebUI.ClickedCh:
				log.Info().Msg("launching web browser")
				openURL(fmt.Sprintf("http://localhost:%d", t.conf.Get().Server.Port))
			case <-mServer.ClickedCh:
				log.Info().Msg("restarting frost")
				t.cancel()
				t.startServices()
			case <-mUI.ClickedCh:
				log.Info().Msg("starting UI")
				go t.startUI()
			case <-mQuit.ClickedCh:
				log.Info().Msg("exiting frost")
				t.cancel()
				systray.Quit()
			}
		}
	}()
}

func (t *Tray) onExit() {
	t.cancel()
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
