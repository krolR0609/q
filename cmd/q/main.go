package main

import (
	"fmt"
	"os"

	"github.com/krolR0609/q/config"
	"github.com/krolR0609/q/internal/app"
	"github.com/krolR0609/q/utils"
)

func main() {
	args := utils.ParseArgs(os.Args[1:])

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Unable to load config:", err)
		return
	}
	app := app.NewApp(cfg)
	app.Run(args)
}
