package log

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"manala/internal/ui"
	"manala/internal/ui/components"
	"testing"
)

type SlogSuite struct {
	suite.Suite
}

func TestSlogSuite(t *testing.T) {
	suite.Run(t, new(SlogSuite))
}

func (s *SlogSuite) TestHandler() {
	outputMock := &ui.OutputMock{}
	outputMock.On("Message", mock.Anything)

	handler := NewSlogHandler(outputMock)
	logger := slog.New(handler)

	logger.Info("message", "foo", "bar")

	outputMock.AssertCalled(s.T(), "Message", &components.Message{
		Type:    components.InfoMessageType,
		Message: "message",
		Attributes: []*components.MessageAttribute{
			{Key: "foo", Value: "bar"},
		},
	})
}
