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

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
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
				"content": fmt.Sprintf("Information about system (use it to make response more relative to the user): %s", systemInfoJson),
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

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d", response.StatusCode)
	}

	var resp OpenAIResponse
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return resp.Choices[0].Message.Content, nil
}
