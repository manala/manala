package log

import (
	"context"
	"log/slog"
	"manala/internal/ui/components"
	"manala/internal/ui/output"
)

func NewSlogHandler(out output.Output) *SlogHandler {
	return &SlogHandler{
		level: slog.LevelInfo,
		out:   out,
	}
}

type SlogHandler struct {
	level  slog.Level
	out    output.Output
	attrs  []slog.Attr
	groups []string
}

func (handler *SlogHandler) LevelDebug() {
	handler.level = slog.LevelDebug
}

func (handler *SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= handler.level
}

func (handler *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := handler.clone()
	clone.attrs = append(clone.attrs, attrs...)

	return clone
}

func (handler *SlogHandler) WithGroup(group string) slog.Handler {
	clone := handler.clone()
	clone.groups = append(clone.groups, group)

	return clone
}

func (handler *SlogHandler) clone() *SlogHandler {
	return &SlogHandler{
		level:  handler.level,
		out:    handler.out,
		attrs:  handler.attrs,
		groups: handler.groups,
	}
}

func (handler *SlogHandler) Handle(_ context.Context, record slog.Record) error {
	message := &components.Message{
		Message: record.Message,
	}

	// Level
	switch record.Level {
	case slog.LevelDebug:
		message.Type = components.DebugMessageType
	case slog.LevelInfo:
		message.Type = components.InfoMessageType
	case slog.LevelWarn:
		message.Type = components.WarnMessageType
	case slog.LevelError:
		message.Type = components.ErrorMessageType
	}

	var attrs []slog.Attr

	// Collect attrs
	attrs = append(attrs, handler.attrs...)
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})

	for _, attr := range attrs {
		message.Attributes = append(message.Attributes, &components.MessageAttribute{
			Key:   attr.Key,
			Value: attr.Value.Any(),
		})
	}

	handler.out.Message(message)

	return nil
}
