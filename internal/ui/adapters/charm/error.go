package charm

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"manala/internal/serrors"
	"manala/internal/ui/components"
)

func (adapter *Adapter) Error(err error) {
	renderer := adapter.errRenderer

	style := messageStyle.New(renderer)

	_, _ = renderer.Output().WriteString(
		style.Render(
			adapter.message(
				adapter.errorMessage(err, renderer),
				renderer,
			),
		) + "\n",
	)
}

func (adapter *Adapter) errorMessage(err error, renderer *lipgloss.Renderer) *components.Message {
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
		message.Details = _err.ErrorDetails(
			renderer.ColorProfile() != termenv.Ascii,
		)
	}

	// Wrapped
	switch _err := err.(type) {
	case interface{ Unwrap() error }:
		message.Messages = append(message.Messages, adapter.errorMessage(_err.Unwrap(), renderer))
	case interface{ Unwrap() []error }:
		for _, __err := range _err.Unwrap() {
			message.Messages = append(message.Messages, adapter.errorMessage(__err, renderer))
		}
	}

	return message
}
