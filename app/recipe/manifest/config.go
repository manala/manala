package manifest

import (
	"github.com/manala/manala/app/sync"
)

type Config struct {
	Description string      `yaml:"description"`
	Icon        string      `yaml:"icon"`
	Template    string      `yaml:"template"`
	Partials    []string    `yaml:"partials"`
	Sync        []sync.Unit `yaml:"sync"`
}
