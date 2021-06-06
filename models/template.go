package models

import (
	"io"
	"manala/template"
)

/***********/
/* Manager */
/***********/

// NewTemplateManager creates a model template manager
func NewTemplateManager(manager template.ManagerInterface, fsManager FsManagerInterface) *templateManager {
	return &templateManager{
		templateManager: manager,
		fsManager:       fsManager,
	}
}

type TemplateManagerInterface interface {
	NewRecipeTemplate(recipe RecipeInterface) (*recipeTemplate, error)
}

type templateManager struct {
	templateManager template.ManagerInterface
	fsManager       FsManagerInterface
}

/*********************/
/* Template - Recipe */
/*********************/

// NewRecipeTemplate creates a recipe template
func (manager *templateManager) NewRecipeTemplate(recipe RecipeInterface) (*recipeTemplate, error) {
	// Fs
	fs := manager.fsManager.NewModelFs(recipe)

	// Recipe template
	tmpl := &recipeTemplate{
		Template: manager.templateManager.NewFsTemplate(fs),
		recipe:   recipe,
	}

	// Include template helpers if any
	helpers := "_helpers.tmpl"
	if _, err := fs.Stat(helpers); err == nil {
		if err := tmpl.ParseFiles(helpers); err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

type recipeTemplate struct {
	*template.Template
	recipe RecipeInterface
}

func (tmpl *recipeTemplate) Execute(writer io.Writer, vars map[string]interface{}) error {
	return tmpl.Template.Execute(writer, map[string]interface{}{
		"Recipe": tmpl.recipe,
		"Vars":   vars,
	})
}
