package main

import "github.com/Alexisdopest/PhoneBridge/internal/app"

func main() {
	// Hide console window logic could be injected here if compiled with -ldflags="-H windowsgui"
	application := app.NewApp()
	application.Run()
}
