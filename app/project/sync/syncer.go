package sync

import (
	"log/slog"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/template"
	"github.com/manala/manala/internal/sync"
)

type Syncer struct {
	syncer         *sync.Syncer
	templateEngine *template.Engine
}

func NewSyncer(log *slog.Logger, templateEngine *template.Engine) *Syncer {
	return &Syncer{
		syncer:         sync.NewSyncer(log),
		templateEngine: templateEngine,
	}
}

func (syncer *Syncer) Sync(project app.Project) error {
	// Template executor
	templateExecutor, err := syncer.templateEngine.Executor(
		project.Vars(),
		project.Recipe(),
		project.Dir(),
	)
	if err != nil {
		return err
	}

	// Loop over project recipe sync units
	for _, unit := range project.Recipe().Sync() {
		if err := syncer.syncer.Sync(
			project.Recipe().Dir(),
			unit.Source(),
			project.Dir(),
			unit.Destination(),
			templateExecutor,
		); err != nil {
			return err
		}
	}

	return nil
}
