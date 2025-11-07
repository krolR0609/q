package app

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

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

	if args.IsChat {
		a.runChatMode()
	} else {
		a.runSingleQuery(args.Prompt)
	}
}

func (a *App) runSingleQuery(prompt string) {
	cancel := utils.ShowLoader()
	_, err := a.ai.Ask(prompt, func() { cancel() })
	if err != nil {
		cancel()
		fmt.Println(err)
		return
	}
}

func (a *App) runChatMode() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Entering chat mode. Type 'exit' to quit.")
	for {
		fmt.Print("q> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "exit" {
			break
		}
		if input == "" {
			continue
		}
		_, err := a.ai.Ask(input, nil)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
