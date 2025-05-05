// Copyright 2025 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package dotazure

import (
	_fs "io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAzdContext(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(string, afero.Fs) error
		opts     []AzdContextOption
		wantPath string
		wantErr  error
	}{
		{
			name:    "no project",
			wantErr: ErrNoProject,
		},
		{
			name: "no environment dir",
			setup: func(root string, fs afero.Fs) error {
				return afero.WriteFile(fs, filepath.Join(root, "src", "azure.yaml"), nil, 0644)
			},
			wantErr: _fs.ErrNotExist,
		},
		{
			name: "no default environment",
			setup: func(root string, fs afero.Fs) error {
				if err := afero.WriteFile(fs, filepath.Join(root, "src", "azure.yaml"), nil, 0644); err != nil {
					return err
				}
				return afero.WriteFile(fs, filepath.Join(root, "src", ".azure", "config.json"), []byte(`{}`), 0644)
			},
			wantErr: ErrNoEnvironmentName,
		},
		{
			name: "default environment",
			setup: func(root string, fs afero.Fs) error {
				if err := afero.WriteFile(fs, filepath.Join(root, "src", "azure.yaml"), nil, 0644); err != nil {
					return err
				}
				if err := afero.WriteFile(fs, filepath.Join(root, "src", ".azure", "config.json"), []byte(`{"defaultEnvironment":"dev"}`), 0644); err != nil {
					return err
				}
				return afero.WriteFile(fs, filepath.Join(root, "src", ".azure", "dev", ".env"), nil, 0644)
			},
			wantPath: filepath.Join("src", ".azure", "dev", ".env"),
		},
		{
			name: "with environment name",
			setup: func(root string, fs afero.Fs) error {
				if err := afero.WriteFile(fs, filepath.Join(root, "src", "azure.yaml"), nil, 0644); err != nil {
					return err
				}
				if err := afero.WriteFile(fs, filepath.Join(root, "src", ".azure", "config.json"), []byte(`{"defaultEnvironment":"dev"}`), 0644); err != nil {
					return err
				}
				return afero.WriteFile(fs, filepath.Join(root, "src", ".azure", "prod", ".env"), nil, 0644)
			},
			opts: []AzdContextOption{
				WithEnvironmentName("prod"),
			},
			wantPath: filepath.Join("src", ".azure", "prod", ".env"),
		},
	}

	// Use the real PWD because filepath.Abs uses it.
	root, err := os.Getwd()
	require.NoError(t, err)
	root, err = filepath.Abs(root)
	require.NoError(t, err)
	pwd := filepath.Join(root, "src", "wd")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a rooted in-memory file system always with the actual PWD as root.
			fs := afero.NewMemMapFs()
			err = fs.MkdirAll(pwd, 0755)
			require.NoError(t, err)
			if tt.setup != nil {
				err := tt.setup(root, fs)
				require.NoError(t, err)
			}
			opts := append([]AzdContextOption{
				// cspell:disable-next-line
				WithFS(fs),
				WithCurrentDirectory(pwd),
			}, tt.opts...)
			context, err := NewAzdContext(opts...)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, filepath.Join(root, tt.wantPath), context.EnvironmentFile())
		})
	}
}

func WithFS(fs afero.Fs) AzdContextOption {
	return func(c *config) error {
		c.fs = fs
		return nil
	}
}
