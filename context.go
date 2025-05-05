// Copyright 2025 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package dotazure

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func init() {
	defaultFs = afero.NewOsFs()
}

var (
	ErrNoProject         = errors.New("no project exists; to create a new project, run `azd init`")
	ErrNoEnvironmentName = errors.New("'.azure/config.json' does not define `defaultEnvironment`")

	defaultFs afero.Fs
)

// AzdContext contains project information for the Azure Development CLI.
type AzdContext struct {
	projectDirectory string
	environmentName  string
}

// NewAzdContext constructs a new AzdContext from the given options.
func NewAzdContext(opts ...AzdContextOption) (*AzdContext, error) {
	c := config{}
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, fmt.Errorf("setting option: %w", err)
		}
	}

	// Get the current working directory.
	var err error
	if c.currentDirectory == "" {
		c.currentDirectory, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("getting current directory: %w", err)
		}
	}

	// Try to find the directory containing azure.yaml.
	searchDir, err := filepath.Abs(c.currentDirectory)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}

	for {
		projectPath := filepath.Join(searchDir, "azure.yaml")
		stat, err := c.FS().Stat(projectPath)
		if os.IsNotExist(err) || (err == nil && stat.IsDir()) {
			parentDir := filepath.Dir(searchDir)
			if parentDir == searchDir {
				return nil, ErrNoProject
			}
			searchDir = parentDir
		} else if err == nil {
			// Found the directory containing azure.yaml.
			break
		} else {
			return nil, fmt.Errorf("searching for project file: %w", err)
		}
	}

	// Get the environment name from .azure/config.json if not set.
	if c.environmentName == "" {
		environmentRoot := filepath.Join(searchDir, ".azure")
		if stat, err := c.FS().Stat(environmentRoot); err != nil || !stat.IsDir() {
			return nil, fmt.Errorf("checking config file: %w", fs.ErrNotExist)
		}
		environmentPath := filepath.Join(environmentRoot, "config.json")
		content, err := afero.ReadFile(c.FS(), environmentPath)
		if err != nil {
			return nil, fmt.Errorf("reading config file: %w", err)
		}

		var config configFile
		if err := json.Unmarshal(content, &config); err != nil {
			return nil, fmt.Errorf("deserializing config file: %w", err)
		}

		if config.DefaultEnvironment == "" {
			return nil, ErrNoEnvironmentName
		}

		c.environmentName = config.DefaultEnvironment
	}

	return &AzdContext{
		projectDirectory: searchDir,
		environmentName:  c.environmentName,
	}, nil
}

// AzdContextOption configures options for constructing an AzdContext.
type AzdContextOption func(*config) error

// WithCurrentDirectory sets a different current directory for NewAzdContext.
func WithCurrentDirectory(path string) AzdContextOption {
	return func(c *config) error {
		if stat, err := c.FS().Stat(path); err != nil || !stat.IsDir() {
			return err
		}
		c.currentDirectory = path
		return nil
	}
}

// WithEnvironmentName sets custom environment name for NewAzdContext.
func WithEnvironmentName(name string) AzdContextOption {
	return func(c *config) error {
		if len(name) == 0 {
			return fmt.Errorf("name cannot be empty")
		}
		c.environmentName = name
		return nil
	}
}

// ProjectDirectory gets the directory containing the azure.yaml project file.
func (c *AzdContext) ProjectDirectory() string {
	return c.projectDirectory
}

// ProjectPath gets the path to the azure.yaml project file.
func (c *AzdContext) ProjectPath() string {
	return filepath.Join(c.ProjectDirectory(), "azure.yaml")
}

// EnvironmentDirectory gets the path to the .azure directory.
func (c *AzdContext) EnvironmentDirectory() string {
	return filepath.Join(c.ProjectDirectory(), ".azure")
}

// EnvironmentName gets the name of the environment.
func (c *AzdContext) EnvironmentName() string {
	return c.environmentName
}

// EnvironmentRoot gets the path to the environment directory under EnvironmentDirectory.
func (c *AzdContext) EnvironmentRoot() string {
	return filepath.Join(c.EnvironmentDirectory(), c.environmentName)
}

// EnvironmentFile gets the path to the .env file under EnvironmentRoot.
func (c *AzdContext) EnvironmentFile() string {
	return filepath.Join(c.EnvironmentRoot(), ".env")
}

type config struct {
	currentDirectory string
	environmentName  string
	fs               afero.Fs
}

func (c *config) FS() afero.Fs {
	if c.fs == nil {
		c.fs = defaultFs
	}
	return c.fs
}

type configFile struct {
	DefaultEnvironment string `json:"defaultEnvironment,omitempty"`
}
