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

type Recipe struct {
	Name       string
	Dir        string
	ConfigFile string
	Config     struct {
		Description string `validate:"required"`
		Sync        []SyncUnit
	}
	Vars map[string]interface{}
}

type SyncUnit struct {
	Source      string
	Destination string
}

// New a recipe
func New(name string) *Recipe {
	return &Recipe{
		Name:       name,
		ConfigFile: configFile,
	}
}

// Load a recipe
func Load(repo *repository.Repository, name string) (*Recipe, error) {
	rec := New(name)

	log.WithField("name", rec.Name).Debug("Loading template...")

	rec.Dir = path.Join(repo.Dir, rec.Name)

	// Load config file
	file, err := os.Open(path.Join(rec.Dir, rec.ConfigFile))
	if err != nil {
		return nil, err
	}

	// Parse
	if err := yaml.NewDecoder(file).Decode(&rec.Vars); err != nil {
		return nil, err
	}

	// Config
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     &rec.Config,
		DecodeHook: StringToSyncUnitHookFunc(),
	})
	if err := decoder.Decode(rec.Vars["manala"]); err != nil {
		return nil, err
	}

	delete(rec.Vars, "manala")

	// Validate
	validate := validator.New()
	if err := validate.Struct(rec); err != nil {
		return nil, err
	}

	return rec, nil
}

// Returns a DecodeHookFunc that converts strings to syncUnit
func StringToSyncUnitHookFunc() mapstructure.DecodeHookFunc {
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
func Walk(repo *repository.Repository, fn walkFunc) error {

	files, err := ioutil.ReadDir(repo.Dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Exclude dot files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.IsDir() {
			rec, err := Load(repo, file.Name())
			if err != nil {
				return err
			}
			fn(rec)
		}
	}

	return nil
}

type walkFunc func(rec *Recipe)
