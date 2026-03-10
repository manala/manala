package sync

import (
	"log/slog"
	"github.com/manala/manala/app"
	"github.com/manala/manala/internal/sync"
)

func NewSyncer(log *slog.Logger) *Syncer {
	return &Syncer{
		syncer: sync.NewSyncer(log),
	}
}

type Syncer struct {
	syncer *sync.Syncer
}

func (syncer *Syncer) Sync(project app.Project) error {
	// Loop over project recipe sync units
	for _, unit := range project.Recipe().Sync() {
		if err := syncer.syncer.Sync(
			project.Recipe().Dir(),
			unit.Source(),
			project.Dir(),
			unit.Destination(),
			project,
		); err != nil {
			return err
		}
	}

	return nil
}
