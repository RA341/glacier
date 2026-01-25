package main

import (
	"log"

	"github.com/ra341/glacier/internal/app"
	"github.com/ra341/glacier/internal/info"
)

func init() {
	app.InitMeta(info.FlavourDocker)
}

func main() {
	file, err := app.LoadUIFromDir("./web")
	if err != nil {
		log.Fatalf("could not load UI from file:%s\nerr:%v", file, err)
		return
	}

	app.NewServer(app.WithUIFS(file))
}
