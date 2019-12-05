package project

import (
	"github.com/apex/log"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

var configFile = ".manala.yaml"

type Project struct {
	Dir        string
	ConfigFile string
	Config     struct {
		Recipe     string `validate:"required"`
		Repository string
	}
	Vars map[string]interface{}
}

func (prj *Project) IsExist() bool {
	info, err := os.Stat(path.Join(prj.Dir, prj.ConfigFile))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Create a project
func New(dir string) *Project {
	return &Project{
		Dir:        dir,
		ConfigFile: configFile,
	}
}

// Load a project
func Load(dir string, repo string) (*Project, error) {
	prj := New(dir)
	prj.Config.Repository = repo

	log.WithField("dir", prj.Dir).Debug("Loading project...")

	// Load config file
	file, err := os.Open(path.Join(dir, prj.ConfigFile))
	if err != nil {
		return nil, err
	}

	// Parse
	if err := yaml.NewDecoder(file).Decode(&prj.Vars); err != nil {
		return nil, err
	}

	// Config
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: &prj.Config,
	})
	if err := decoder.Decode(prj.Vars["manala"]); err != nil {
		return nil, err
	}

	delete(prj.Vars, "manala")

	// Validate
	validate := validator.New()
	if err := validate.Struct(prj); err != nil {
		return nil, err
	}

	return prj, nil
}
