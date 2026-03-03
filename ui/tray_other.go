//go:build !windows

package main

import "context"

func startTray(_ context.Context) {
	// System tray is only supported on Windows.
}
