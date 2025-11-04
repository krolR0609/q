package app

import (
	"fmt"
	"net/http"

	"github.com/krolR0609/q/config"
	"github.com/krolR0609/q/internal/api/openai"
	"github.com/krolR0609/q/internal/httpclient"
	"github.com/krolR0609/q/internal/services/systeminfo"
	"github.com/krolR0609/q/utils"
)

type App struct {
	config     *config.Config
	httpClient *http.Client
	ai         *openai.Provider
	sysInfo    systeminfo.SystemInfoProvider
}

func NewApp(config *config.Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) InitServices() {
	a.httpClient = httpclient.NewDefaultClient(a.config)
	a.sysInfo = systeminfo.NewBasicSystemInfoProvider()
	a.ai = openai.NewOpenAiProvider(a.config, a.httpClient, a.sysInfo)
}

func (a *App) Run(args utils.Args) {
	a.InitServices()

	cancel := utils.ShowLoader()
	result, err := a.ai.Ask(args.Prompt)
	cancel()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
