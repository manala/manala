package project

import (
	"fmt"
	"github.com/apex/log"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

var configFile = ".manala.yaml"

// Create a project
func New(dir string) Interface {
	return &project{
		dir: dir,
	}
}

type Interface interface {
	GetDir() string
	GetConfigFile() string
	GetConfig() Config
	IsExist() bool
	GetVars() map[string]interface{}
	Load(cfg Config) error
}

type project struct {
	dir    string
	config Config
	vars   map[string]interface{}
}

type Config struct {
	Recipe     string `validate:"required"`
	Repository string
}

func (prj *project) GetDir() string {
	return prj.dir
}

func (prj *project) GetConfigFile() string {
	return path.Join(prj.dir, configFile)
}

func (prj *project) GetConfig() Config {
	return prj.config
}

func (prj *project) IsExist() bool {
	info, err := os.Stat(prj.GetConfigFile())
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (prj *project) GetVars() map[string]interface{} {
	return prj.vars
}

// Load project
func (prj *project) Load(cfg Config) error {
	// Project exist ?
	if !prj.IsExist() {
		return fmt.Errorf("project not found")
	}

	prj.config = cfg

	log.WithField("dir", prj.dir).Debug("Loading project...")

	// Load config file
	file, err := os.Open(prj.GetConfigFile())
	if err != nil {
		return err
	}

	// Parse vars
	if err := yaml.NewDecoder(file).Decode(&prj.vars); err != nil {
		return err
	}

	// Map config
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: &prj.config,
	})
	if err := decoder.Decode(prj.vars["manala"]); err != nil {
		return err
	}

	delete(prj.vars, "manala")

	// Validate
	validate := validator.New()
	if err := validate.Struct(prj.config); err != nil {
		return err
	}

	return nil
}
