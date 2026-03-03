//go:build windows

package main

import (
	"context"
	_ "embed"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/appicon.png
var trayIcon []byte

func startTray(ctx context.Context) {
	systray.Run(
		func() {
			systray.SetIcon(trayIcon)
			systray.SetTooltip("OpenVINO Desktop")

			mShow := systray.AddMenuItem("Show", "Open the window")
			systray.AddSeparator()
			mQuit := systray.AddMenuItem("Quit", "Exit OpenVINO Desk")

			for {
				select {
				case <-mShow.ClickedCh:
					runtime.WindowShow(ctx)
				case <-mQuit.ClickedCh:
					runtime.Quit(ctx)
					return
				}
			}
		},
		func() {},
	)
}
