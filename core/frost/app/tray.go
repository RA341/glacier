package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"runtime"
	"sync"
	"sync/atomic"

	"fyne.io/systray"
)

type Tray struct {
	ctx    context.Context
	cancel context.CancelFunc
	subFS  fs.FS

	wg            sync.WaitGroup
	uiRunning     atomic.Bool
	serverRunning atomic.Bool
}

func NewTray(subFs fs.FS) {
	t := Tray{
		subFS: subFs,
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
	// hold until functions cleans up
	t.wg.Wait()

	// then reset context
	ctx, cancel := context.WithCancel(context.Background())
	t.ctx = ctx
	t.cancel = cancel

	runtime.GC()

	t.wg.Go(t.startServer)
	t.wg.Go(t.startUI)
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

func (t *Tray) loadIcon() []byte {
	open, err := t.subFS.Open("favicon.png")
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

func (t *Tray) onExit() {
	t.cancel()
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

	NewServer(WithUIFS(t.subFS), WithCtx(t.ctx))
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

	err := NewUI(t.ctx, "ui/ui")
	if err != nil {
		if errors.Is(t.ctx.Err(), context.Canceled) {
			fmt.Println("Process stopped by user")
			return
		}

		ShowErr(fmt.Sprintf("Failed to start UI: %v", err))
	}
}
