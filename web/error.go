package web

import (
	"context"
	"manala/internal/ui/components"
	"net/http"
)

func (server *Server) error(ctx context.Context, response http.ResponseWriter, err error) {
	server.out.Error(err)

	message := components.MessageFromError(err, true)

	if err := server.template(ctx, response, "error.gohtml", map[string]any{
		"Message": message,
	}); err != nil {
		server.out.Error(err)
	}
}
