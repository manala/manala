package url

import (
	"log/slog"
	"github.com/manala/manala/app"
	"github.com/manala/manala/app/repository"
)

func NewProcessorLoaderHandler(log *slog.Logger, processor *Processor) *ProcessorLoaderHandler {
	return &ProcessorLoaderHandler{
		log:       log.With("handler", "url.processor"),
		processor: processor,
	}
}

type ProcessorLoaderHandler struct {
	log       *slog.Logger
	processor *Processor
}

func (handler *ProcessorLoaderHandler) Handle(query *repository.LoaderQuery, chain repository.LoaderHandlerChain) (app.Repository, error) {
	handler.log.Debug("handle repository url", "url", query.URL)

	var err error

	// Process query url
	query.URL, err = handler.processor.Process(query.URL)
	if err != nil {
		return nil, err
	}

	// Chain
	return chain.Next(query)
}
