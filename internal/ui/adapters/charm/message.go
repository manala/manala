package charm

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"manala/internal/ui/components"
	"strings"
)

func (adapter *Adapter) Message(message *components.Message) {
	renderer := adapter.errRenderer

	style := messageStyle.New(renderer)

	_, _ = renderer.Output().WriteString(
		style.Render(
			adapter.message(message, renderer),
		) + "\n",
	)
}

func (adapter *Adapter) message(message *components.Message, renderer *lipgloss.Renderer) string {
	// Empty message
	if message.Message == "" {
		return ""
	}

	var symbolStyle, attributeKeyStyle *style

	switch message.Type {
	case components.DebugMessageType:
		symbolStyle = debugSymbolStyle.New(renderer)
		attributeKeyStyle = debugStyle.New(renderer)
	case components.InfoMessageType:
		symbolStyle = infoSymbolStyle.New(renderer)
		attributeKeyStyle = infoStyle.New(renderer)
	case components.WarnMessageType:
		symbolStyle = warnSymbolStyle.New(renderer)
		attributeKeyStyle = warnStyle.New(renderer)
	case components.ErrorMessageType:
		symbolStyle = errorSymbolStyle.New(renderer)
		attributeKeyStyle = errorStyle.New(renderer)
	}

	messageStyle := messageMessageStyle.New(renderer)
	attributesStyle := messageAttributesStyle.New(renderer)
	attributeValueStyle := messageAttributeValueStyle.New(renderer)
	detailsStyle := messageDetailsStyle.New(renderer)

	// Message
	str := messageStyle.Render(
		message.Message,
	)

	// Attributes
	attributes := make([]string, len(message.Attributes))
	for i := range message.Attributes {
		attributes[i] = attributeKeyStyle.Render(message.Attributes[i].Key) + "=" +
			attributeValueStyle.Render(fmt.Sprintf("%v", message.Attributes[i].Value))
	}

	if len(attributes) > 0 {
		str = lipgloss.JoinVertical(lipgloss.Left,
			str,
			attributesStyle.Render(
				strings.Join(attributes, " "),
			),
		)
	}

	// Details
	if message.Details != "" {
		str = lipgloss.JoinVertical(lipgloss.Left,
			str,
			detailsStyle.Render(
				message.Details,
			),
		)
	}

	for _, _message := range message.Messages {
		str = lipgloss.JoinVertical(lipgloss.Left,
			str,
			adapter.message(_message, renderer),
		)
	}

	// Symbol
	return lipgloss.JoinHorizontal(lipgloss.Top,
		symbolStyle.Render(""),
		str,
	)
}
