package name

import (
	"log/slog"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe"
)

func NewProcessorLoaderHandler(log *slog.Logger, processor *Processor) *ProcessorLoaderHandler {
	return &ProcessorLoaderHandler{
		log:       log.With("handler", "name.processor"),
		processor: processor,
	}
}

type ProcessorLoaderHandler struct {
	log       *slog.Logger
	processor *Processor
}

func (handler *ProcessorLoaderHandler) Handle(query *recipe.LoaderQuery, chain recipe.LoaderHandlerChain) (app.Recipe, error) {
	handler.log.Debug("handle recipe name", "name", query.Name)

	// Process query name
	query.Name = handler.processor.Process(query.Name)

	// Chain
	return chain.Next(query)
}
