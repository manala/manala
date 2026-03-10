package log_test

import (
	"log/slog"
	"github.com/manala/manala/internal/ui"
	"github.com/manala/manala/internal/ui/components"
	"github.com/manala/manala/internal/ui/log"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SlogSuite struct {
	suite.Suite
}

func TestSlogSuite(t *testing.T) {
	suite.Run(t, new(SlogSuite))
}

func (s *SlogSuite) TestHandler() {
	s.Run("Default", func() {
		outputMock := &ui.OutputMock{}
		outputMock.On("Message", mock.Anything)

		handler := log.NewSlogHandler(outputMock)
		log := slog.New(handler)

		log.Info("info", "foo", "bar")
		log.Debug("debug", "foo", "bar")

		outputMock.AssertNumberOfCalls(s.T(), "Message", 1)
		outputMock.AssertCalled(s.T(), "Message", &components.Message{
			Type:    components.InfoMessageType,
			Message: "info",
			Attributes: []*components.MessageAttribute{
				{Key: "foo", Value: "bar"},
			},
		})
	})

	s.Run("WithDebug", func() {
		outputMock := &ui.OutputMock{}
		outputMock.On("Message", mock.Anything)

		handler := log.NewSlogHandler(outputMock,
			log.WithSlogHandlerDebug(true),
		)
		log := slog.New(handler)

		log.Info("info", "foo", "bar")
		log.Debug("debug", "foo", "bar")

		outputMock.AssertExpectations(s.T())

		outputMock.AssertNumberOfCalls(s.T(), "Message", 2)
		outputMock.AssertCalled(s.T(), "Message", &components.Message{
			Type:    components.InfoMessageType,
			Message: "info",
			Attributes: []*components.MessageAttribute{
				{Key: "foo", Value: "bar"},
			},
		})
		outputMock.AssertCalled(s.T(), "Message", &components.Message{
			Type:    components.DebugMessageType,
			Message: "debug",
			Attributes: []*components.MessageAttribute{
				{Key: "foo", Value: "bar"},
			},
		})
	})
}
