package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"

	"fyne.io/systray"
	"github.com/ra341/glacier/shared/api"
)

type Tray struct {
	ctx    context.Context
	cancel context.CancelFunc

	conf TrayConfig

	wg            sync.WaitGroup
	uiRunning     atomic.Bool
	serverRunning atomic.Bool
}

func NewTray(opts ...Opt) {
	var conf TrayConfig
	for _, opt := range opts {
		opt(&conf)
	}

	t := Tray{
		conf: conf,
	}
	t.Start()
}

func (t *Tray) Start() {
	all := t.loadIcon()

	t.startServices()
	systray.Run(
		func() {
			t.onReady(all)
		},
		t.onExit,
	)
}

func (t *Tray) startServices() {
	// hold until functions clean up
	t.wg.Wait()

	// then reset context
	ctx, cancel := context.WithCancel(context.Background())
	t.ctx = ctx
	t.cancel = cancel

	t.wg.Go(t.startServer)
	t.wg.Go(t.startUI)
}

func (t *Tray) startServer() {
	if t.serverRunning.Load() {
		fmt.Println("Server is already running")
		return
	}

	fmt.Println("Starting server...")

	defer func() {
		fmt.Println("Server stopped")
		t.serverRunning.Store(false)
	}()
	t.serverRunning.Store(true)

	NewServer(api.WithServerBase(&t.conf.serverBase))
}

func (t *Tray) startUI() {
	if t.uiRunning.Load() {
		fmt.Println("UI is running")
		return
	}

	fmt.Println("Starting UI")
	defer func() {
		t.uiRunning.Store(false)
	}()
	t.uiRunning.Store(true)

	if t.conf.disableUI {
		fmt.Println("UI is disabled in config")
		return
	}

	exePath := "ui/ui"
	if runtime.GOOS == "windows" {
		exePath += ".exe"
	}

	err := NewUI(t.ctx, exePath)
	if err != nil {
		if errors.Is(t.ctx.Err(), context.Canceled) {
			fmt.Println("Process stopped by user")
			return
		}

		ShowErr(fmt.Sprintf("Failed to start UI: %v", err))
	}
}

func (t *Tray) onReady(all []byte) {
	systray.SetIcon(all)
	systray.SetTitle("Frost")
	systray.SetTooltip("Frost")

	mUI := systray.AddMenuItem("Open UI", "Start the UI")
	mServer := systray.AddMenuItem("Restart", "Restart the app")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		for {
			select {
			case <-mServer.ClickedCh:
				t.cancel()
				t.startServices()
			case <-mUI.ClickedCh:
				go t.startUI()
			case <-mQuit.ClickedCh:
				t.cancel()
				systray.Quit()
			}
		}
	}()
}

func (t *Tray) onExit() {
	t.cancel()
}

func (t *Tray) loadIcon() []byte {
	uifs := t.conf.serverBase.UIFS
	if uifs == nil {
		return nil
	}

	open, err := uifs.Open("favicon.png")
	if err != nil {
		ShowErr("Could not open favicon.svg")
		os.Exit(1)
	}

	all, err := io.ReadAll(open)
	if err != nil {
		ShowErr("Could not read favicon.svg bytes")
		os.Exit(1)
	}

	return all
}
