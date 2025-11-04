package httpclient

import (
	"net/http"
	"time"

	"github.com/krolR0609/q/config"
)

func NewDefaultClient(config *config.Config) *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        3,
		MaxIdleConnsPerHost: 3,
		MaxConnsPerHost:     3,
		IdleConnTimeout:     30 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   2 * time.Minute, // can be quite long
	}

	return client
}
