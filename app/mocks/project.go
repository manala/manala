package mocks

import (
	"github.com/stretchr/testify/mock"
	"io"
	"manala/app/interfaces"
	"manala/internal/template"
	"manala/internal/validation"
)

type ProjectMock struct {
	mock.Mock
}

func (rec *ProjectMock) Dir() string {
	args := rec.Called()
	return args.String(0)
}

func (rec *ProjectMock) Recipe() interfaces.Recipe {
	args := rec.Called()
	return args.Get(0).(interfaces.Recipe)
}

func (rec *ProjectMock) Vars() map[string]interface{} {
	args := rec.Called()
	return args.Get(0).(map[string]interface{})
}

func (rec *ProjectMock) Template() *template.Template {
	args := rec.Called()
	return args.Get(0).(*template.Template)
}

/************/
/* Manifest */
/************/

type ProjectManifestMock struct {
	mock.Mock
}

func (man *ProjectManifestMock) Recipe() string {
	args := man.Called()
	return args.String(0)
}

func (man *ProjectManifestMock) Repository() string {
	args := man.Called()
	return args.String(0)
}

func (man *ProjectManifestMock) Vars() map[string]interface{} {
	args := man.Called()
	return args.Get(0).(map[string]interface{})
}

func (man *ProjectManifestMock) ReadFrom(reader io.Reader) error {
	args := man.Called(reader)
	return args.Error(0)
}

func (man *ProjectManifestMock) ValidationResultErrorDecorator() validation.ResultErrorDecorator {
	args := man.Called()
	return args.Get(0).(validation.ResultErrorDecorator)
}
