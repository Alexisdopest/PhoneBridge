package tray

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Alexisdopest/PhoneBridge/internal/autostart"
	"github.com/getlantern/systray"
)

type TrayManager struct {
	onQuit func()
}

func NewTrayManager(onQuit func()) *TrayManager {
	return &TrayManager{
		onQuit: onQuit,
	}
}

func (t *TrayManager) Start() {
	systray.Run(t.onReady, t.onExit)
}

func (t *TrayManager) Stop() {
	systray.Quit()
}

func (t *TrayManager) onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("PhoneBridge")
	systray.SetTooltip("PhoneBridge - LAN Collaboration")

	mOpenDir := systray.AddMenuItem("打开接收文件夹", "Open the folder where files are saved")
	
	mAutoStart := systray.AddMenuItemCheckbox("开机启动", "Start PhoneBridge with Windows", autostart.IsEnabled())
	
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "Quit PhoneBridge")

	go func() {
		for {
			select {
			case <-mOpenDir.ClickedCh:
				homeDir, _ := os.UserHomeDir()
				destDir := filepath.Join(homeDir, "Downloads", "PhoneBridge")
				os.MkdirAll(destDir, 0755)
				exec.Command("explorer", destDir).Start()

			case <-mAutoStart.ClickedCh:
				if mAutoStart.Checked() {
					if err := autostart.Disable(); err == nil {
						mAutoStart.Uncheck()
					}
				} else {
					if err := autostart.Enable(); err == nil {
						mAutoStart.Check()
					}
				}

			case <-mQuit.ClickedCh:
				t.onQuit()
				// systray.Quit() is intentionally called here or in app lifecycle
				return
			}
		}
	}()
}

func (t *TrayManager) onExit() {
	// Cleanup if needed
}
