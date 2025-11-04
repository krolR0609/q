package systeminfo

import (
	"runtime"
	"time"
)

type BasicSystemInfoProvider struct{}

func NewBasicSystemInfoProvider() *BasicSystemInfoProvider {
	return &BasicSystemInfoProvider{}
}

func (p BasicSystemInfoProvider) GetSystemInfo() map[string]string {
	info := make(map[string]string)

	info["GOOS"] = runtime.GOOS
	info["GOARCH"] = runtime.GOARCH
	info["CurrentTime"] = time.Now().Format(time.RFC3339)

	return info
}
