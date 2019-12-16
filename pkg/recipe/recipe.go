package recipe

import (
	"github.com/apex/log"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"manala/pkg/repository"
	"os"
	"path"
	"reflect"
	"strings"
)

var configFile = ".manala.yaml"

// Create a recipe
func New(name string) Interface {
	return &recipe{
		name: name,
	}
}

type Interface interface {
	GetName() string
	GetDir() string
	GetConfig() Config
	GetVars() map[string]interface{}
	Load(repo repository.Interface) error
}

type recipe struct {
	name   string
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
	return rec.name
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

func (rec *recipe) GetVars() map[string]interface{} {
	return rec.vars
}

// Load recipe
func (rec *recipe) Load(repo repository.Interface) error {
	log.WithField("name", rec.name).Debug("Loading recipe...")

	rec.dir = path.Join(repo.GetDir(), rec.name)

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
	if err := validate.Struct(rec); err != nil {
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

// walk into repository recipes
func Walk(repo repository.Interface, fn walkFunc) error {

	files, err := ioutil.ReadDir(repo.GetDir())
	if err != nil {
		return err
	}

	for _, file := range files {
		// Exclude dot files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.IsDir() {
			rec := New(file.Name())
			if err := rec.Load(repo); err != nil {
				return err
			}
			fn(rec)
		}
	}

	return nil
}

type walkFunc func(rec Interface)
