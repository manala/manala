package mocks

import (
	"github.com/manala/manala/app"

	"github.com/stretchr/testify/mock"
)

// Project mock a Project.
type Project struct {
	mock.Mock
}

func (p *Project) Dir() string {
	args := p.Called()

	return args.String(0)
}

func (p *Project) Recipe() app.Recipe {
	args := p.Called()

	return args.Get(0).(app.Recipe)
}

func (p *Project) Vars() map[string]any {
	args := p.Called()

	return args.Get(0).(map[string]any)
}

func (p *Project) Watches() ([]string, error) {
	args := p.Called()

	return args.Get(0).([]string), args.Error(1)
}
