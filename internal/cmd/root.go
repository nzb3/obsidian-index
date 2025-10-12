package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nzb3/obsidian-index/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "obsidian-index",
	Short: "A CLI tool for indexing Obsidian vaults",
	Long: `obsidian-index is a powerful CLI tool that creates comprehensive
indexes for your Obsidian vault by generating markdown files with links
to all entries in each directory, processing from leaves to root.`,
}

func Execute() {
	rootCmd.ParseFlags(os.Args[1:])

	versionFlag, err := rootCmd.Flags().GetBool("version")
	if err != nil {
		slog.Error("failed to parse version flag", "error", err)
		fmt.Fprintf(os.Stderr, "Error parsing version flag: %v\n", err)
		os.Exit(1)
	}

	if versionFlag {
		fmt.Println(version.String())
		return
	}

	if err := rootCmd.Execute(); err != nil {
		slog.Error("error executing command", "error", err)
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")
}
