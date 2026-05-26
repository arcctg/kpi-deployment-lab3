package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `host: 127.0.0.1
port: 5000
database:
  host: localhost
  port: 5432
  user: appuser
  password: secret
  dbname: notes
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	cfg, err := LoadConfig(path)
	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.Equal(t, 5000, cfg.Port)
	assert.Equal(t, "appuser", cfg.Database.User)
	assert.Equal(t, "secret", cfg.Database.Password)
}

func TestLoadConfigEnvOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `host: 127.0.0.1
port: 5000
database:
  host: localhost
  port: 5432
  user: appuser
  password: secret
  dbname: notes
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	t.Setenv("APP_HOST", "0.0.0.0")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("DB_USER", "override")
	t.Setenv("DB_PASSWORD", "newpass")

	cfg, err := LoadConfig(path)
	require.NoError(t, err)
	assert.Equal(t, "0.0.0.0", cfg.Host)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "override", cfg.Database.User)
	assert.Equal(t, "newpass", cfg.Database.Password)
}
