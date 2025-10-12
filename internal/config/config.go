package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	vaultDir    string
	verbose     bool
	dryRun      bool
	backup      bool
	excludeDirs []string
}

func New() *Config {
	return &Config{
		vaultDir:    "",
		verbose:     false,
		dryRun:      false,
		backup:      false,
		excludeDirs: []string{},
	}
}

func NewWithOptions(vaultDir string, verbose bool) *Config {
	return &Config{
		vaultDir:    vaultDir,
		verbose:     verbose,
		dryRun:      false,
		backup:      false,
		excludeDirs: []string{},
	}
}

func NewWithAllOptions(vaultDir string, verbose, dryRun, backup bool, excludeDirs []string) *Config {
	return &Config{
		vaultDir:    vaultDir,
		verbose:     verbose,
		dryRun:      dryRun,
		backup:      backup,
		excludeDirs: excludeDirs,
	}
}

func (c *Config) GetVaultDir() string {
	return c.vaultDir
}

func (c *Config) IsVerbose() bool {
	return c.verbose
}

func (c *Config) IsDryRun() bool {
	return c.dryRun
}

func (c *Config) IsBackup() bool {
	return c.backup
}

func (c *Config) GetExcludeDirs() []string {
	return c.excludeDirs
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.vaultDir == "" {
		return errors.New("vault directory is required")
	}

	// Check if vault directory exists and is accessible
	info, err := os.Stat(c.vaultDir)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("vault directory does not exist: " + c.vaultDir)
		}
		return errors.New("cannot access vault directory: " + err.Error())
	}

	if !info.IsDir() {
		return errors.New("vault path is not a directory: " + c.vaultDir)
	}

	// Check if vault directory is absolute path
	if !filepath.IsAbs(c.vaultDir) {
		return errors.New("vault directory must be an absolute path: " + c.vaultDir)
	}

	// Validate exclude directories
	for _, dir := range c.excludeDirs {
		if strings.TrimSpace(dir) == "" {
			return errors.New("exclude directory cannot be empty")
		}
	}

	return nil
}
