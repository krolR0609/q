package api

type BaseEventResponder struct {
	onFirstChunk func()
}

func NewEventResponder(onFirstChunk func()) *BaseEventResponder {
	return &BaseEventResponder{
		onFirstChunk: onFirstChunk,
	}
}

func (e *BaseEventResponder) OnFirstChunk() {
	e.onFirstChunk()
}
