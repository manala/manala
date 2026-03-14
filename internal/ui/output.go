package ui

import (
	"github.com/manala/manala/internal/ui/components"

	"github.com/stretchr/testify/mock"
)

type Output interface {
	Message(message *components.Message)
	Error(err error)
}

type OutputMock struct {
	mock.Mock
}

func (mock *OutputMock) Message(message *components.Message) {
	mock.Called(message)
}

func (mock *OutputMock) Error(err error) {
	mock.Called(err)
}
