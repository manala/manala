package name

import (
	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/internal/log"
)

type ProcessorLoaderHandler struct {
	log       *log.Log
	processor *Processor
}

func NewProcessorLoaderHandler(log *log.Log, processor *Processor) *ProcessorLoaderHandler {
	return &ProcessorLoaderHandler{
		log:       log,
		processor: processor,
	}
}

func (handler *ProcessorLoaderHandler) Handle(query *recipe.LoaderQuery, chain recipe.LoaderHandlerChain) (app.Recipe, error) {
	handler.log.Debug("handle recipe name", "handler", "name.processor", "name", query.Name)

	// Process query name
	query.Name = handler.processor.Process(query.Name)

	// Chain
	return chain.Next(query)
}
