package lipgloss

import (
	"manala/internal/errors/serrors"
	"manala/internal/ui/components"
)

func (output *Output) Error(err error) {
	output.Message(
		output.errorMessage(err),
	)
}

func (output *Output) errorMessage(err error) *components.Message {
	message := &components.Message{
		Type:    components.ErrorMessageType,
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
				message.Attributes = append(message.Attributes, &components.MessageAttribute{
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
		message.Details = _err.ErrorDetails(output.errAnsi())
	}

	// Wrapped
	switch _err := err.(type) {
	case interface{ Unwrap() error }:
		message.Messages = append(message.Messages, output.errorMessage(_err.Unwrap()))
	case interface{ Unwrap() []error }:
		for _, __err := range _err.Unwrap() {
			message.Messages = append(message.Messages, output.errorMessage(__err))
		}
	}

	return message
}
