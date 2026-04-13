package charm

import (
	"fmt"
	"strings"

	"github.com/manala/manala/internal/ui/components"

	"charm.land/lipgloss/v2"
)

func (adapter *Adapter) Message(message *components.Message) {
	style := messageStyle

	_, _ = lipgloss.Fprintln(adapter.err,
		style.Render(
			adapter.message(message),
		),
	)
}

func (adapter *Adapter) message(message *components.Message) string {
	// Empty message
	if message.Message == "" {
		return ""
	}

	var symbolStyle, attributeKeyStyle lipgloss.Style

	switch message.Type {
	case components.DebugMessageType:
		symbolStyle = debugSymbolStyle
		attributeKeyStyle = debugStyle
	case components.InfoMessageType:
		symbolStyle = infoSymbolStyle
		attributeKeyStyle = infoStyle
	case components.WarnMessageType:
		symbolStyle = warnSymbolStyle
		attributeKeyStyle = warnStyle
	case components.ErrorMessageType:
		symbolStyle = errorSymbolStyle
		attributeKeyStyle = errorStyle
	}

	messageStyle := messageMessageStyle
	attributesStyle := messageAttributesStyle
	attributeValueStyle := messageAttributeValueStyle
	dumpStyle := messageDumpStyle

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
		str = lipgloss.JoinHorizontal(lipgloss.Top,
			str,
			attributesStyle.Render(
				strings.Join(attributes, " "),
			),
		)
	}

	// Dump
	if message.Dump != "" {
		str = lipgloss.JoinVertical(lipgloss.Left,
			str,
			dumpStyle.Render(
				message.Dump,
			),
		)
	}

	for _, _message := range message.Messages {
		str = lipgloss.JoinVertical(lipgloss.Left,
			str,
			adapter.message(_message),
		)
	}

	// Symbol
	return lipgloss.JoinHorizontal(lipgloss.Top,
		symbolStyle.Render(""),
		str,
	)
}
