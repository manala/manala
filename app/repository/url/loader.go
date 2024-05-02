package url

import (
	"log/slog"
	"manala/app"
	"manala/app/repository"
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
	handler.log.Debug("handle repository url", "url", query.Url)

	var err error

	// Process query url
	query.Url, err = handler.processor.Process(query.Url)
	if err != nil {
		return nil, err
	}

	// Chain
	return chain.Next(query)
}
