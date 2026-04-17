package url

import (
	"github.com/manala/manala/app"
	"github.com/manala/manala/app/repository"
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

func (handler *ProcessorLoaderHandler) Handle(query *repository.LoaderQuery, chain repository.LoaderHandlerChain) (app.Repository, error) {
	handler.log.Debug("handle repository url", "handler", "url.processor", "url", query.URL)

	var err error

	// Process query url
	query.URL, err = handler.processor.Process(query.URL)
	if err != nil {
		return nil, err
	}

	// Chain
	return chain.Next(query)
}
