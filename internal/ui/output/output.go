package output

import "manala/internal/ui/components"

type Output interface {
	Message(message *components.Message)
	Error(err error)
	Table(table *components.Table)
}
