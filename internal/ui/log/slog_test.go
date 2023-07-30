package log

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"manala/internal/ui/components"
	"manala/internal/ui/output"
	"testing"
)

type SlogSuite struct {
	suite.Suite
}

func TestSlogSuite(t *testing.T) {
	suite.Run(t, new(SlogSuite))
}

func (s *SlogSuite) TestHandler() {
	outputMock := &output.Mock{}
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
