package systeminfo

type SystemInfoProvider interface {
	GetSystemInfo() map[string]string
}
