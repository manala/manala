package log

import (
	"context"
	"log/slog"
	"manala/internal/ui"
	"manala/internal/ui/components"
)

func NewSlogHandler(out ui.Output, opts ...SlogHandlerOption) *SlogHandler {
	handler := &SlogHandler{
		out: out,
	}

	for _, opt := range opts {
		opt(handler)
	}

	return handler
}

type SlogHandler struct {
	debug  bool
	out    ui.Output
	attrs  []slog.Attr
	groups []string
}

func (handler *SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	if handler.debug {
		return level >= slog.LevelDebug
	}

	return level >= slog.LevelInfo
}

func (handler *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := *handler
	clone.attrs = append(clone.attrs, attrs...)

	return &clone
}

func (handler *SlogHandler) WithGroup(group string) slog.Handler {
	clone := *handler
	clone.groups = append(clone.groups, group)

	return &clone
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

type SlogHandlerOption func(handler *SlogHandler)

func WithSlogHandlerDebug(debug bool) SlogHandlerOption {
	return func(handler *SlogHandler) {
		handler.debug = debug
	}
}
