//go:build windows

package main

import (
	"os"

	"golang.org/x/sys/windows/registry"
)

const (
	runKey  = `Software\Microsoft\Windows\CurrentVersion\Run`
	appName = "OpenVINO"
)

// SetStartup adds or removes the app from the Windows startup registry key.
func (a *App) SetStartup(enable bool) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	if enable {
		exe, err := os.Executable()
		if err != nil {
			return err
		}
		return k.SetStringValue(appName, `"`+exe+`"`)
	}
	return k.DeleteValue(appName)
}

// GetStartupEnabled reports whether the app is registered to run at startup.
func (a *App) GetStartupEnabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	_, _, err = k.GetStringValue(appName)
	return err == nil
}
