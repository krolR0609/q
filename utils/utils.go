package utils

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Args struct {
	Prompt string
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

	prompt := strings.Join(nonFlags, " ")
	if p, ok := dict["-r"]; ok {
		prompt = strings.TrimSpace(p)
	}

	parsedArgs := Args{
		Prompt: strings.TrimSpace(prompt),
	}

	return parsedArgs
}

func ShowLoader() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		loader := []string{"|", "/", "-", "\\"}
		i := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Print("\r\033[K")
				return
			default:
				// Keep printing loading animation
				fmt.Print("\r", loader[i%len(loader)])
				i++
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	return cancel
}
