package indexator

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Indexator struct {
	vaultPath   string
	dryRun      bool
	backup      bool
	excludeDirs []string
}

func NewIndexator(vaultPath string) *Indexator {
	return &Indexator{
		vaultPath:   vaultPath,
		dryRun:      false,
		backup:      false,
		excludeDirs: []string{},
	}
}

func NewIndexatorWithOptions(vaultPath string, dryRun, backup bool, excludeDirs []string) *Indexator {
	return &Indexator{
		vaultPath:   vaultPath,
		dryRun:      dryRun,
		backup:      backup,
		excludeDirs: excludeDirs,
	}
}

// Start begins the indexing process, starting from leaves and moving to root
func (idx *Indexator) Start() error {
	directories, err := idx.CollectDirectories()
	if err != nil {
		slog.Error("failed to collect directories", "error", err)
		return fmt.Errorf("failed to collect directories: %w", err)
	}

	sort.Slice(directories, func(i, j int) bool {
		depthI := strings.Count(directories[i], string(filepath.Separator))
		depthJ := strings.Count(directories[j], string(filepath.Separator))
		return depthI > depthJ
	})

	for _, dir := range directories {
		if err := idx.indexDirectory(dir); err != nil {
			slog.Error("failed to index directory", "directory", dir, "error", err)
			return fmt.Errorf("failed to index directory %s: %w", dir, err)
		}
	}

	return nil
}

func (idx *Indexator) CollectDirectories() ([]string, error) {
	var directories []string

	err := fs.WalkDir(os.DirFS(idx.vaultPath), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				slog.Warn("permission denied, skipping", "path", path, "error", err)
				return nil
			}
			slog.Error("error walking directory", "path", path, "error", err)
			return err
		}

		if d.IsDir() && strings.HasPrefix(d.Name(), ".") && path != "." {
			return filepath.SkipDir
		}

		if d.IsDir() {
			// Check if directory should be excluded
			if idx.shouldExcludeDirectory(path) {
				return filepath.SkipDir
			}
			directories = append(directories, path)
		}

		return nil
	})

	return directories, err
}

func (idx *Indexator) indexDirectory(dirPath string) error {
	fullPath := filepath.Join(idx.vaultPath, dirPath)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		slog.Error("failed to read directory", "directory", dirPath, "error", err)
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var links []string

	for _, entry := range entries {
		entryPath := filepath.Join(fullPath, entry.Name())

		if idx.isIndexFile(entryPath, fullPath) {
			continue
		}

		if entry.IsDir() {
			indexFileName := entry.Name() + ".md"
			indexPath := filepath.Join(entryPath, indexFileName)

			if _, err := os.Stat(indexPath); err == nil {
				relPath := idx.getRelativePath(indexPath)
				links = append(links, fmt.Sprintf("[[%s]]", relPath))
			}
		} else {
			relPath := idx.getRelativePath(entryPath)
			links = append(links, fmt.Sprintf("[[%s]]", relPath))
		}
	}

	if len(links) == 0 {
		return nil
	}

	// Check if index file already exists
	dirName := filepath.Base(fullPath)
	if dirName == "." || fullPath == idx.vaultPath {
		dirName = "index"
	}
	indexFileName := dirName + ".md"
	indexFilePath := filepath.Join(fullPath, indexFileName)

	if _, err := os.Stat(indexFilePath); err == nil {
		// Index file already exists, skip creation
		return nil
	}

	return idx.createIndexFile(fullPath, links)
}

func (idx *Indexator) createIndexFile(dirPath string, links []string) error {
	dirName := filepath.Base(dirPath)
	if dirName == "." || dirPath == idx.vaultPath {
		dirName = "index"
	}

	indexFileName := dirName + ".md"
	indexFilePath := filepath.Join(dirPath, indexFileName)

	content := strings.Join(links, "\n") + "\n"

	// Handle dry run mode
	if idx.dryRun {
		slog.Info("DRY RUN: Would create index", "file", indexFilePath, "entries", len(links))
		return nil
	}

	// Handle backup if file exists
	if idx.backup {
		if err := idx.backupExistingFile(indexFilePath); err != nil {
			slog.Warn("failed to backup existing file", "file", indexFilePath, "error", err)
		}
	}

	// Use atomic file operation to prevent race conditions
	return idx.writeFileAtomic(indexFilePath, []byte(content))
}

// writeFileAtomic writes content to a file atomically to prevent race conditions
func (idx *Indexator) writeFileAtomic(filePath string, content []byte) error {
	// Create temporary file in the same directory
	tempFile := filePath + ".tmp"

	// Write to temporary file first
	err := os.WriteFile(tempFile, content, 0644)
	if err != nil {
		slog.Error("failed to write temporary file", "file", tempFile, "error", err)
		return fmt.Errorf("failed to write temporary file %s: %w", tempFile, err)
	}

	// Atomic rename operation
	err = os.Rename(tempFile, filePath)
	if err != nil {
		// Clean up temporary file on failure
		os.Remove(tempFile)
		slog.Error("failed to rename temporary file", "temp", tempFile, "target", filePath, "error", err)
		return fmt.Errorf("failed to rename temporary file %s to %s: %w", tempFile, filePath, err)
	}

	slog.Info("Created index", "file", filePath, "entries", len(strings.Split(strings.TrimSpace(string(content)), "\n")))
	return nil
}

func (idx *Indexator) getRelativePath(absolutePath string) string {
	relPath, err := filepath.Rel(idx.vaultPath, absolutePath)
	if err != nil {
		return absolutePath
	}

	return strings.ReplaceAll(relPath, string(filepath.Separator), "/")
}

func (idx *Indexator) isIndexFile(filePath, dirPath string) bool {
	fileName := filepath.Base(filePath)
	dirName := filepath.Base(dirPath)

	if dirPath == idx.vaultPath {
		dirName = "index"
	}

	expectedIndexName := dirName + ".md"
	return fileName == expectedIndexName
}

// shouldExcludeDirectory checks if a directory should be excluded from indexing
func (idx *Indexator) shouldExcludeDirectory(dirPath string) bool {
	for _, excludeDir := range idx.excludeDirs {
		// Check if the directory path contains the exclude pattern
		if strings.Contains(dirPath, excludeDir) {
			return true
		}
		// Check if the directory name matches the exclude pattern
		dirName := filepath.Base(dirPath)
		if dirName == excludeDir {
			return true
		}
	}
	return false
}

// backupExistingFile creates a backup of an existing file with timestamp
func (idx *Indexator) backupExistingFile(filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // No file to backup
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filePath + ".backup_" + timestamp

	// Copy file to backup location
	err := os.Rename(filePath, backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup %s: %w", backupPath, err)
	}

	slog.Info("Created backup", "original", filePath, "backup", backupPath)
	return nil
}
