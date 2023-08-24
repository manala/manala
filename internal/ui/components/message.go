package components

type MessageType int

const (
	DebugMessageType MessageType = iota - 1
	InfoMessageType
	WarnMessageType
	ErrorMessageType
)

type Message struct {
	Type       MessageType
	Message    string
	Attributes []*MessageAttribute
	Details    string
	Messages   []*Message
}

type MessageAttribute struct {
	Key   string
	Value any
}
