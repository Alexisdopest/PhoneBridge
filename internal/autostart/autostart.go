package autostart

import (
	"log"
	"os"

	"golang.org/x/sys/windows/registry"
)

const appName = "PhoneBridge"

// Enable adds the current executable to the Windows startup registry
func Enable() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	err = key.SetStringValue(appName, exePath)
	if err != nil {
		return err
	}

	log.Println("Successfully enabled Windows auto-start")
	return nil
}

// Disable removes the application from the Windows startup registry
func Disable() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	err = key.DeleteValue(appName)
	if err != nil {
		// Ignore if it doesn't exist
		if err == registry.ErrNotExist {
			return nil
		}
		return err
	}

	log.Println("Successfully disabled Windows auto-start")
	return nil
}

// IsEnabled checks if auto-start is currently enabled
func IsEnabled() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	val, _, err := key.GetStringValue(appName)
	if err != nil {
		return false
	}

	exePath, _ := os.Executable()
	return val == exePath
}
