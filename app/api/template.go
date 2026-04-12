package api

import "github.com/manala/manala/app/template"

func (api *API) NewTemplateEngine() *template.Engine {
	return template.NewEngine()
}
