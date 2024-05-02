package name

import (
	"cmp"
	"log/slog"
	"slices"
)

func NewProcessor(log *slog.Logger) *Processor {
	return &Processor{
		log: log,
	}
}

type Processor struct {
	log     *slog.Logger
	entries []processorEntry
}

func (processor *Processor) Add(name string, weight int) {
	processor.entries = append(processor.entries, processorEntry{name: name, weight: weight})
}

func (processor *Processor) Process(name string) string {
	entries := slices.Clone(processor.entries)

	if name != "" {
		entries = append(entries, processorEntry{name: name, weight: 0})
	}

	// Sort entries by weight
	slices.SortFunc(entries, func(a, b processorEntry) int {
		return cmp.Compare(b.weight, a.weight)
	})

	for _, entry := range entries {
		processor.log.Debug("process recipe name",
			"name", entry.name,
			"weight", entry.weight,
		)

		name = entry.name

		if name != "" {
			break
		}
	}

	return name
}

type processorEntry struct {
	name   string
	weight int
}
