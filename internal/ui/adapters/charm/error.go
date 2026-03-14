package charm

import (
	"os"

	"github.com/manala/manala/internal/ui/components"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
)

func (adapter *Adapter) Error(err error) {
	style := messageStyle

	profile := colorprofile.Detect(adapter.err, os.Environ())

	_, _ = lipgloss.Fprintln(adapter.err,
		style.Render(
			adapter.message(
				components.MessageFromError(
					err,
					profile >= colorprofile.ANSI,
				),
			),
		),
	)
}
