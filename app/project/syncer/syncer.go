package syncer

import (
	"log/slog"
	"manala/app"
	"manala/internal/syncer"
)

func New(log *slog.Logger) *Syncer {
	return &Syncer{
		syncer: syncer.New(log),
	}
}

type Syncer struct {
	syncer *syncer.Syncer
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
