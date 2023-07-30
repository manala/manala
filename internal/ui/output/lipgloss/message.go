package lipgloss

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"manala/internal/ui/components"
	"strings"
)

var messageBullets = map[components.MessageType]string{
	components.DebugMessageType: "•",
	components.InfoMessageType:  "•",
	components.WarnMessageType:  "•",
	components.ErrorMessageType: "⨯",
}

func (output *Output) Message(message *components.Message) {
	output.message(message, 0)
}

func (output *Output) message(message *components.Message, indentation int) {
	// Empty message
	if message.Message == "" {
		return
	}

	var style lipgloss.Style

	switch message.Type {
	case components.DebugMessageType:
		style = output.errStyle().Foreground(lipgloss.Color("15")).Bold(true)
	case components.InfoMessageType:
		style = output.errStyle().Foreground(lipgloss.Color("12")).Bold(true)
	case components.WarnMessageType:
		style = output.errStyle().Foreground(lipgloss.Color("11")).Bold(true)
	case components.ErrorMessageType:
		style = output.errStyle().Foreground(lipgloss.Color("9")).Bold(true)
	}

	// Message
	str := output.errStyle().
		MarginLeft(1).
		Width(33).
		Render(message.Message)

	// Attributes
	var attributesBuilder strings.Builder
	for _, attribute := range message.Attributes {
		attributesBuilder.WriteString(fmt.Sprintf(" %s=%v", style.Render(attribute.Key), attribute.Value))
	}
	attributes := attributesBuilder.String()

	// Attributes
	if attributes != "" {
		str = lipgloss.JoinHorizontal(lipgloss.Top, str,
			output.errStyle().
				MarginLeft(1).
				Render(attributes),
		)
	}

	// Details
	if message.Details != "" {
		str = lipgloss.JoinVertical(lipgloss.Left,
			str,
			output.errStyle().
				MarginLeft(1).
				Render("\n"+message.Details),
		)
	}

	// Bullet
	str = lipgloss.JoinHorizontal(lipgloss.Top,
		style.Copy().
			PaddingLeft(2+(indentation*2)).
			Render(messageBullets[message.Type]),
		str,
	)

	// Trim
	for _, line := range strings.Split(str, "\n") {
		output.writeErrString(strings.TrimRight(line, " ") + "\n")
	}

	for _, _message := range message.Messages {
		output.message(_message, indentation+1)
	}
}
