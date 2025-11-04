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
	current := ""
	dict[current] = current
	for _, arg := range args {
		if arg[0] == '-' {
			dict[arg] = ""
			current = arg
		} else {
			dict[current] += strings.TrimSpace(arg) + " "
		}
	}

	prompt, ok := dict["-r"]
	if !ok {
		prompt, _ = dict[""]
	}

	model := Args{
		Prompt: strings.TrimSpace(prompt),
	}

	return model
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
