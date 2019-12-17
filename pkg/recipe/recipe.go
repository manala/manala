package recipe

import (
	"fmt"
	"github.com/apex/log"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"manala/pkg/clean"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

var configFile = ".manala.yaml"

// Create a recipe
func New(dir string) Interface {
	return &recipe{
		dir: dir,
	}
}

type Interface interface {
	GetName() string
	GetDir() string
	GetConfigFile() string
	GetConfig() Config
	IsExist() bool
	GetVars() map[string]interface{}
	Load(cfg Config) error
}

type recipe struct {
	dir    string
	config Config
	vars   map[string]interface{}
}

type Config struct {
	Description string `validate:"required"`
	Sync        []SyncUnit
}

type SyncUnit struct {
	Source      string
	Destination string
}

func (rec *recipe) GetName() string {
	return filepath.Base(rec.dir)
}

func (rec *recipe) GetDir() string {
	return rec.dir
}

func (rec *recipe) GetConfigFile() string {
	return path.Join(rec.dir, configFile)
}

func (rec *recipe) GetConfig() Config {
	return rec.config
}

func (rec *recipe) IsExist() bool {
	info, err := os.Stat(rec.GetConfigFile())
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (rec *recipe) GetVars() map[string]interface{} {
	return rec.vars
}

// Load recipe
func (rec *recipe) Load(cfg Config) error {
	// Recipe exist ?
	if !rec.IsExist() {
		return fmt.Errorf("recipe not found")
	}

	rec.config = cfg

	log.WithField("name", rec.GetName()).Debug("Loading recipe...")

	// Load config file
	cfgFile, err := os.Open(path.Join(rec.GetConfigFile()))
	if err != nil {
		return err
	}

	// Parse config
	var cfgMap map[string]interface{}
	if err := yaml.NewDecoder(cfgFile).Decode(&cfgMap); err != nil {
		return err
	}

	// See: https://github.com/go-yaml/yaml/issues/139
	cfgMap = clean.YamlStringMap(cfgMap)

	// Map config
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     &rec.config,
		DecodeHook: stringToSyncUnitHookFunc(),
	})
	if err := decoder.Decode(cfgMap["manala"]); err != nil {
		return err
	}

	delete(cfgMap, "manala")
	rec.vars = cfgMap

	// Validate
	validate := validator.New()
	if err := validate.Struct(rec.config); err != nil {
		return err
	}

	return nil
}

// Returns a DecodeHookFunc that converts strings to syncUnit
func stringToSyncUnitHookFunc() mapstructure.DecodeHookFunc {
	return func(rf reflect.Type, rt reflect.Type, data interface{}) (interface{}, error) {
		if rf.Kind() != reflect.String {
			return data, nil
		}
		if rt != reflect.TypeOf(SyncUnit{}) {
			return data, nil
		}

		src := data.(string)
		dst := src

		// Separate source / destination
		u := strings.Split(src, " ")
		if len(u) > 1 {
			src = u[0]
			dst = u[1]
		}

		return SyncUnit{
			Source:      src,
			Destination: dst,
		}, nil
	}
}
