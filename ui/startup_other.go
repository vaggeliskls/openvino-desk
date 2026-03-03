//go:build !windows

package main

import "errors"

func (a *App) SetStartup(_ bool) error {
	return errors.New("startup registration is only supported on Windows")
}

func (a *App) GetStartupEnabled() bool {
	return false
}
