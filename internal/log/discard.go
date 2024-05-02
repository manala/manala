package log

import (
	"context"
	"log/slog"
)

var Discard = slog.New(DiscardHandler{})

type DiscardHandler struct{}

func (DiscardHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (DiscardHandler) Handle(context.Context, slog.Record) error { return nil }
func (d DiscardHandler) WithAttrs([]slog.Attr) slog.Handler      { return d }
func (d DiscardHandler) WithGroup(string) slog.Handler           { return d }
