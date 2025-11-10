package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/nzb3/obsidian-index/internal/app"
	"github.com/nzb3/obsidian-index/internal/config"
	"github.com/spf13/cobra"
)

var (
	vaultDir    string
	verbose     bool
	dryRun      bool
	backup      bool
	excludeDirs []string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize indexing for an Obsidian vault",
	Long: `Initialize and start the indexing process for an Obsidian vault.
This command will recursively index all directories starting from the
deepest level (leaves) and working up to the root directory.

Each directory will get an index file named after the directory containing
markdown links to all files and subdirectories within it.`,
	Example: `  obsidian-index init
  obsidian-index init --dir /path/to/obsidian/vault
  obsidian-index init -d ~/Documents/MyVault --verbose`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&vaultDir, "dir", "d", "", "path to the Obsidian vault directory (default: current directory)")

	initCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	initCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be done without creating files")
	initCmd.Flags().BoolVar(&backup, "backup", false, "create backup of existing index files")
	initCmd.Flags().StringSliceVar(&excludeDirs, "exclude", []string{}, "directories to exclude from indexing")
}

func runInit(cmd *cobra.Command, args []string) error {
	dir := vaultDir
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			slog.Error("failed to get current working directory", "error", err)
			return fmt.Errorf("failed to get current working directory: %w", err)
		}
	}

	absPath, err := filepath.Abs(dir)
	if err != nil {
		slog.Error("failed to get absolute path", "directory", dir, "error", err)
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	cfg := config.NewWithAllOptions(absPath, verbose, dryRun, backup, excludeDirs)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		slog.Error("configuration validation failed", "error", err)
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	if verbose {
		fmt.Printf("Starting indexation of vault: %s\n", absPath)
		if dryRun {
			fmt.Println("🔍 DRY RUN MODE - No files will be created")
		}
		if backup {
			fmt.Println("💾 BACKUP MODE - Existing index files will be backed up")
		}
		if len(excludeDirs) > 0 {
			fmt.Printf("🚫 Excluding directories: %v\n", excludeDirs)
		}
	}

	application := app.New(cfg)

	if err := application.Run(); err != nil {
		slog.Error("indexation failed", "vault", absPath, "error", err)
		return fmt.Errorf("indexation failed: %w", err)
	}

	if dryRun {
		fmt.Printf("🔍 Dry run completed for vault: %s\n", absPath)
	} else {
		fmt.Printf("✅ Successfully indexed vault: %s\n", absPath)
	}
	return nil
}
