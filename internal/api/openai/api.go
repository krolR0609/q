package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/krolR0609/q/config"
	"github.com/krolR0609/q/internal/services/systeminfo"
)

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

func (p *Provider) Ask(prompt string) (string, error) {
	baseUrl, err := url.Parse(fmt.Sprintf("%s/chat/completions", p.config.BaseUrl))
	if err != nil {
		return "", err
	}

	systemInfo := p.systemInfo.GetSystemInfo()
	systemInfoJson, err := json.Marshal(systemInfo)
	if err != nil {
		return "", err
	}

	reqBody := map[string]interface{}{
		"model":       p.config.Model,
		"max_tokens":  32000,
		"temperature": 0.55,
		"top_p":       1,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are \"q\" CLI tool which allows to query different AI tools and return the response as user wished. Return info in short form because it will be displayed in the terminal! If user asking about terminal tool or specific command return only command in the response!",
			},
			{
				"role":    "system",
				"content": fmt.Sprintf("Infomation about system (use it to make response more relative to the user): %s", systemInfoJson),
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	response, err := p.httpClient.Post(baseUrl.String(), "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	rDict := map[string]interface{}{}
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&rDict)
	choices, ok := rDict["choices"].([]interface{})
	if !ok {
		panic("choices is not a slice")
	}

	// Get first choice
	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		panic("first choice is not a map")
	}

	// Get message
	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		panic("message is not a map")
	}

	// Get content
	content, ok := message["content"].(string)
	if !ok {
		panic("content is not a string")
	}

	return content, nil
}
