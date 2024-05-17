package web

import (
	"context"
	"embed"
	"html/template"
	"manala/app"
	"manala/internal/serrors"
	"manala/internal/ui/components"
	"net/http"

	"github.com/Masterminds/sprig/v3"
	"github.com/robert-nix/ansihtml"
)

//go:embed templates
var templates embed.FS

func (server *Server) mustTemplate(ctx context.Context, response http.ResponseWriter, name string, data map[string]any) {
	if err := server.template(ctx, response, name, data); err != nil {
		server.error(ctx, response, serrors.NewTemplate(err))
	}
}

func (server *Server) template(ctx context.Context, response http.ResponseWriter, name string, data map[string]any) error {
	// Funcs
	funcs := sprig.FuncMap()

	funcs["messageTypeClass"] = func(messageType components.MessageType) string {
		switch messageType {
		case components.DebugMessageType:
			return "secondary"
		case components.WarnMessageType:
			return "warning"
		case components.ErrorMessageType:
			return "danger"
		default:
			return "primary"
		}
	}

	funcs["ansiToHTML"] = func(ansi string) template.HTML {
		return template.HTML(
			ansihtml.ConvertToHTML([]byte(ansi)),
		)
	}

	tmpl, err := template.New("").
		Funcs(funcs).
		ParseFS(templates, "templates/layout.gohtml", "templates/"+name)
	if err != nil {
		return serrors.NewTemplate(err)
	}

	// Data
	data["RepositoryUrl"], _ = app.RepositoryURL(ctx)
	data["RepositoryRef"], _ = app.RepositoryRef(ctx)

	if err = tmpl.ExecuteTemplate(response, "layout.gohtml", data); err != nil {
		return serrors.NewTemplate(err)
	}

	return nil
}
