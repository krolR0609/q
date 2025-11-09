package api

type LLMWriter interface {
	Ask(messages []map[string]string, evt LLMEventResponder) (string, error)
}

type LLMEventResponder interface {
	OnFirstChunk()
}
