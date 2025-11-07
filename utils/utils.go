package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

type Args struct {
	Prompt string
	IsChat bool
}

func ParseArgs(args []string) Args {
	dict := make(map[string]string)
	var nonFlags []string
	currentFlag := ""
	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' {
			currentFlag = arg
			dict[currentFlag] = ""
		} else {
			if currentFlag == "" {
				nonFlags = append(nonFlags, arg)
			} else {
				dict[currentFlag] += arg + " "
			}
		}
	}

	isChat := false
	prompt := strings.Join(nonFlags, " ")
	if len(nonFlags) > 0 && nonFlags[0] == "chat" {
		isChat = true
		// For chat mode, remove "chat" from prompt if more args
		if len(nonFlags) > 1 {
			prompt = strings.Join(nonFlags[1:], " ")
		} else {
			prompt = ""
		}
	}
	if p, ok := dict["-r"]; ok {
		prompt = strings.TrimSpace(p)
	}

	parsedArgs := Args{
		Prompt: strings.TrimSpace(prompt),
		IsChat: isChat,
	}

	return parsedArgs
}

func ShowLoader() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		loader := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Print("\r\033[K")
				os.Stdout.Sync()
				return
			default:
				// Keep printing loading animation
				fmt.Print("\r", loader[i%len(loader)])
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return cancel
}
