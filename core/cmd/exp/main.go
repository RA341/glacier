package main

import (
	"log"

	"github.com/wailsapp/wails/lib/renderer/webview"
)

func main() {
	err := webview.Open(
		"Test",
		"http://localhost:6699",
		1200,
		720,
		true,
	)
	if err != nil {
		log.Fatal(err)
		return
	}
}
