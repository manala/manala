package validation

import (
	"errors"
	"regexp"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

var (
	// Git Repo.
	gitRepoRegex  = regexp.MustCompile(`((git|ssh|http(s)?)|(git@[\w.]+))(:(//)?)([\w.@:/\-~]+)(\.git)(/)?`)
	GitRepoFormat = &jsonschema.Format{
		Name: "git-repo",
		Validate: func(v any) error {
			s, ok := v.(string)
			if !ok {
				return errors.New("not a string")
			}
			if !gitRepoRegex.MatchString(s) {
				return errors.New("invalid format")
			}
			return nil
		},
	}
	// File Path.
	filePathRegex  = regexp.MustCompile(`^/[a-z0-9-_/]*$`)
	FilePathFormat = &jsonschema.Format{
		Name: "file-path",
		Validate: func(v any) error {
			s, ok := v.(string)
			if !ok {
				return errors.New("not a string")
			}
			if !filePathRegex.MatchString(s) {
				return errors.New("invalid format")
			}
			return nil
		},
	}
	// Domain.
	domainRegex  = regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.)+[a-zA-Z]{2,}$`)
	DomainFormat = &jsonschema.Format{
		Name: "domain",
		Validate: func(v any) error {
			s, ok := v.(string)
			if !ok {
				return errors.New("not a string")
			}
			if !domainRegex.MatchString(s) {
				return errors.New("invalid format")
			}
			return nil
		},
	}
)
