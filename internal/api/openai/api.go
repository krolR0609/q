package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/krolR0609/q/config"
	"github.com/krolR0609/q/internal/api"
	"github.com/krolR0609/q/internal/services/systeminfo"
)

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
	Delta   Delta   `json:"delta"`
}

type Message struct {
	Content string `json:"content"`
}

type Delta struct {
	Content string `json:"content"`
}

type Provider struct {
	httpClient *http.Client
	config     *config.Config
	systemInfo systeminfo.SystemInfoProvider
}

func NewOpenAiProvider(
	config *config.Config,
	httpClient *http.Client,
	systemInfo systeminfo.SystemInfoProvider,
) *Provider {
	return &Provider{
		config:     config,
		httpClient: httpClient,
		systemInfo: systemInfo,
	}
}

func (p *Provider) Ask(messages []map[string]string, evt api.LLMEventResponder) (string, error) {
	baseUrl, err := url.Parse(fmt.Sprintf("%s/chat/completions", p.config.BaseUrl))
	if err != nil {
		return "", err
	}

	systemInfo := p.systemInfo.GetSystemInfo()
	systemInfoJson, err := json.Marshal(systemInfo)
	if err != nil {
		return "", err
	}

	systemMessages := []map[string]string{
		{
			"role":    "system",
			"content": "You are \"q\" CLI tool which allows to query different AI tools and return the response as user wished. Return info in short form because it will be displayed in the terminal! If user asking about terminal tool or specific command return only command in the response!",
		},
		{
			"role":    "system",
			"content": fmt.Sprintf("Information about system (use it to make response more relative to the user): %s", systemInfoJson),
		},
	}

	reqBody := map[string]interface{}{
		"model":       p.config.Model,
		"max_tokens":  32000,
		"temperature": 0.55,
		"top_p":       1,
		"stream":      true,
		"messages":    append(systemMessages, messages...),
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", baseUrl.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	if p.config.Options != nil && p.config.Options.ApiKey != "" {
		request.Header.Add("Authorization", p.config.Options.ApiKey)
	}

	response, err := p.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("sending request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d", response.StatusCode)
	}

	scanner := bufio.NewScanner(response.Body)
	var fullResponse strings.Builder
	var allData strings.Builder
	hasStream := false
	called := false
	firstContent := true
	for scanner.Scan() {
		line := scanner.Text()
		allData.WriteString(line + "\n")
		if strings.HasPrefix(line, "data: ") {
			hasStream = true
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}
			var chunk OpenAIResponse
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				return "", err
			}
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				if !called && evt != nil {
					evt.OnFirstChunk()
					called = true
				}
				content := chunk.Choices[0].Delta.Content
				if firstContent {
					fmt.Print("\n" + content)
					firstContent = false
				} else {
					fmt.Print(content)
				}
				os.Stdout.Sync()
				fullResponse.WriteString(content)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if !hasStream {
		// Fallback to non-streamed response
		var resp OpenAIResponse
		if err := json.Unmarshal([]byte(strings.TrimSpace(allData.String())), &resp); err != nil {
			return "", err
		}
		if len(resp.Choices) > 0 {
			if evt != nil {
				evt.OnFirstChunk()
			}
			content := resp.Choices[0].Message.Content
			fmt.Print("\n" + content)
			os.Stdout.Sync()
			fullResponse.WriteString(content)
		}
	}
	fmt.Println() // Add newline after response
	os.Stdout.Sync()
	return fullResponse.String(), nil
}
