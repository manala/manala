package models

import (
	"github.com/stretchr/testify/suite"
	"manala/fs"
	"manala/template"
	"testing"
)

/*********/
/* Suite */
/*********/

type TemplateTestSuite struct {
	suite.Suite
	manager TemplateManagerInterface
}

func TestTemplateTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(TemplateTestSuite))
}

func (s *TemplateTestSuite) SetupTest() {
	fsManager := fs.NewManager()
	templateManager := template.NewManager()
	modelFsManager := NewFsManager(fsManager)
	s.manager = NewTemplateManager(templateManager, modelFsManager)
}

/*********/
/* Tests */
/*********/

func (s *TemplateTestSuite) Test() {
	repository := NewRepository("foo", "foo", false)
	recipe := NewRecipe("foo", "foo", "foo", "foo", repository, nil, nil, nil, nil)
	tmpl, err := s.manager.NewRecipeTemplate(recipe)
	s.NoError(err)
	s.Implements((*template.Interface)(nil), tmpl)
}
