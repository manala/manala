package url

import (
	"cmp"
	"fmt"
	"log/slog"
	"manala/internal/serrors"
	netURL "net/url"
	"slices"
	"strings"

	"dario.cat/mergo"
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

func (processor *Processor) Add(url string, weight int) {
	processor.entries = append(processor.entries, processorEntry{url: url, weight: weight})
}

func (processor *Processor) AddQuery(key string, value string, weight int) {
	query := make(netURL.Values)
	query.Set(key, value)
	processor.entries = append(processor.entries, processorEntry{query: query, weight: weight})
}

func (processor *Processor) Process(url string) (string, error) {
	entries := slices.Clone(processor.entries)

	if url != "" {
		entries = append(entries, processorEntry{url: url, weight: 0})
	}

	// Sort entries by weight
	slices.SortFunc(entries, func(a, b processorEntry) int {
		return cmp.Compare(b.weight, a.weight)
	})

	var query netURL.Values

	for _, entry := range entries {
		// Split url query parts
		var entryQuery string
		entry.url, entryQuery, _ = strings.Cut(entry.url, "?")

		if entryQuery != "" {
			values, err := netURL.ParseQuery(entryQuery)
			if err != nil {
				return "", serrors.New("unable to process repository query").
					WithArguments("query", entryQuery).
					WithErrors(err)
			}
			_ = mergo.Merge(&entry.query, values)
		}

		processor.log.Debug("process repository url",
			"url", entry.url,
			"query", entry.query.Encode(),
			"weight", entry.weight,
		)

		url = entry.url
		_ = mergo.Merge(&query, entry.query)

		if url != "" {
			break
		}
	}

	if url != "" && query != nil {
		url = fmt.Sprintf("%s?%s", url, query.Encode())
	}

	return url, nil
}

type processorEntry struct {
	url    string
	query  netURL.Values
	weight int
}
