package components

import (
	"manala/internal/serrors"
)

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

func MessageFromError(err error, ansi bool) *Message {
	message := &Message{
		Type:    ErrorMessageType,
		Message: err.Error(),
	}

	// Arguments
	if _err, ok := err.(serrors.ErrorArguments); ok {
		arguments := _err.ErrorArguments()
		for len(arguments) > 0 {
			switch key := arguments[0].(type) {
			case string:
				if len(arguments) == 1 {
					arguments = nil
				}
				message.Attributes = append(message.Attributes, &MessageAttribute{
					Key:   key,
					Value: arguments[1],
				})
				arguments = arguments[2:]
			default:
				arguments = arguments[1:]
			}
		}
	}

	// Details
	if _err, ok := err.(serrors.ErrorDetails); ok {
		message.Details = _err.ErrorDetails(ansi)
	}

	// Wrapped
	switch _err := err.(type) {
	case interface{ Unwrap() error }:
		message.Messages = append(message.Messages, MessageFromError(_err.Unwrap(), ansi))
	case interface{ Unwrap() []error }:
		for _, __err := range _err.Unwrap() {
			message.Messages = append(message.Messages, MessageFromError(__err, ansi))
		}
	}

	return message
}
