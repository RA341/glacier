package main

import (
	"log"

	"github.com/ra341/glacier/internal/app"
	"github.com/ra341/glacier/internal/info"
	"github.com/ra341/glacier/shared/api"
)

func init() {
	app.InitMeta(info.FlavourDocker)
}

func main() {
	file, err := api.LoadUIFromDir("./web")
	if err != nil {
		log.Fatalf("could not load UI from file:%s\nerr:%v", file, err)
		return
	}

	app.NewServer(api.WithUIFS(file))
}
