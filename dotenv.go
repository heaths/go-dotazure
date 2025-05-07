// Copyright 2025 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package dotazure

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/joho/godotenv"
)

// Load loads environment variables for an Azure Developer CLI project.
//
// Locates the `.env` file from the default environment name if an Azure Developer CLI project was already provisioned.
// Returns true if a `.env` file was found and loaded successfully; otherwise, returns an error.
func Load(opts ...LoadOption) (bool, error) {
	l := loader{}

	var err error
	for _, opt := range opts {
		if err = opt(&l); err != nil {
			return false, fmt.Errorf("setting option: %w", err)
		}
	}

	if l.context == nil {
		if l.context, err = NewAzdContext(); errors.Is(err, fs.ErrNotExist) {
			return false, nil
		} else if err != nil {
			return false, fmt.Errorf("getting default context: %w", err)
		}
	}

	path := l.context.EnvironmentFile()
	if l.replace {
		err = godotenv.Overload(path)
	} else {
		err = godotenv.Load(path)
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("loading environment variables: %w", err)
	}
	return true, nil
}

// LoadOption configures options for loading environment variables.
type LoadOption func(*loader) error

// WithContext sets the AzdContext to use for discovery.
func WithContext(context *AzdContext) LoadOption {
	return func(l *loader) error {
		l.context = context
		return nil
	}
}

// WithReplace sets whether to replace environment variables that were already set.
func WithReplace(replace bool) LoadOption {
	return func(l *loader) error {
		l.replace = replace
		return nil
	}
}

type loader struct {
	context *AzdContext
	replace bool
}
