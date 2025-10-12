package indexator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewIndexator(t *testing.T) {
	vaultPath := "/test/vault"
	indexator := NewIndexator(vaultPath)

	if indexator == nil {
		t.Fatal("NewIndexator returned nil")
	}

	if indexator.vaultPath != vaultPath {
		t.Errorf("Expected vaultPath %s, got %s", vaultPath, indexator.vaultPath)
	}
}

func TestIndexator_getRelativePath(t *testing.T) {
	tests := []struct {
		name         string
		vaultPath    string
		absolutePath string
		expected     string
	}{
		{
			name:         "normal relative path",
			vaultPath:    "/test/vault",
			absolutePath: "/test/vault/notes/file.md",
			expected:     "notes/file.md",
		},
		{
			name:         "same directory",
			vaultPath:    "/test/vault",
			absolutePath: "/test/vault/file.md",
			expected:     "file.md",
		},
		{
			name:         "subdirectory with spaces",
			vaultPath:    "/test/vault",
			absolutePath: "/test/vault/my notes/important file.md",
			expected:     "my notes/important file.md",
		},
		{
			name:         "error case - returns relative path",
			vaultPath:    "/test/vault",
			absolutePath: "/completely/different/path/file.md",
			expected:     "../../completely/different/path/file.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := &Indexator{vaultPath: tt.vaultPath}
			result := idx.getRelativePath(tt.absolutePath)
			if result != tt.expected {
				t.Errorf("getRelativePath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIndexator_isIndexFile(t *testing.T) {
	tests := []struct {
		name      string
		vaultPath string
		filePath  string
		dirPath   string
		expected  bool
	}{
		{
			name:      "is index file in subdirectory",
			vaultPath: "/test/vault",
			filePath:  "/test/vault/notes/notes.md",
			dirPath:   "/test/vault/notes",
			expected:  true,
		},
		{
			name:      "is index file in root",
			vaultPath: "/test/vault",
			filePath:  "/test/vault/index.md",
			dirPath:   "/test/vault",
			expected:  true,
		},
		{
			name:      "not index file",
			vaultPath: "/test/vault",
			filePath:  "/test/vault/notes/other.md",
			dirPath:   "/test/vault/notes",
			expected:  false,
		},
		{
			name:      "not index file in root",
			vaultPath: "/test/vault",
			filePath:  "/test/vault/readme.md",
			dirPath:   "/test/vault",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := &Indexator{vaultPath: tt.vaultPath}
			result := idx.isIndexFile(tt.filePath, tt.dirPath)
			if result != tt.expected {
				t.Errorf("isIndexFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIndexator_Start_Integration(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Create test directory structure
	testDirs := []string{
		"notes",
		"notes/subfolder",
		"documents",
		"documents/archive",
	}

	for _, dir := range testDirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	// Create some test files
	testFiles := []string{
		"notes/file1.md",
		"notes/file2.md",
		"notes/subfolder/file3.md",
		"documents/doc1.md",
		"documents/archive/old.md",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		err := os.WriteFile(filePath, []byte("# Test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create indexator and run indexing
	indexator := NewIndexator(tempDir)
	err := indexator.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Verify that index files were created for directories with content
	expectedIndexFiles := []string{
		"notes/notes.md",               // notes directory index
		"notes/subfolder/subfolder.md", // subfolder index
		"documents/documents.md",       // documents index
		"documents/archive/archive.md", // archive index
	}

	for _, indexFile := range expectedIndexFiles {
		fullPath := filepath.Join(tempDir, indexFile)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected index file %s was not created", indexFile)
		}
	}

	// Verify content of notes index file
	notesIndexPath := filepath.Join(tempDir, "notes/notes.md")
	content, err := os.ReadFile(notesIndexPath)
	if err != nil {
		t.Fatalf("Failed to read notes index file: %v", err)
	}

	contentStr := string(content)
	expectedLinks := []string{
		"[[notes/file1.md]]",
		"[[notes/file2.md]]",
		"[[notes/subfolder/subfolder.md]]",
	}

	for _, link := range expectedLinks {
		if !strings.Contains(contentStr, link) {
			t.Errorf("Notes index file should contain link %s", link)
		}
	}
}

func TestIndexator_Start_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	indexator := NewIndexator(tempDir)
	err := indexator.Start()
	if err != nil {
		t.Fatalf("Start() failed on empty directory: %v", err)
	}

	// Should not create any index files for empty directory
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp directory: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected empty directory, but found %d entries", len(entries))
	}
}

func TestIndexator_Start_WithHiddenDirectories(t *testing.T) {
	tempDir := t.TempDir()

	// Create hidden directory
	hiddenDir := filepath.Join(tempDir, ".hidden")
	err := os.MkdirAll(hiddenDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create hidden directory: %v", err)
	}

	// Create file in hidden directory
	hiddenFile := filepath.Join(hiddenDir, "secret.md")
	err = os.WriteFile(hiddenFile, []byte("# Secret"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file in hidden directory: %v", err)
	}

	// Create normal directory with files
	normalDir := filepath.Join(tempDir, "normal")
	err = os.MkdirAll(normalDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create normal directory: %v", err)
	}

	normalFile := filepath.Join(normalDir, "public.md")
	err = os.WriteFile(normalFile, []byte("# Public"), 0644)
	if err != nil {
		t.Fatalf("Failed to create normal file: %v", err)
	}

	indexator := NewIndexator(tempDir)
	err = indexator.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Verify that hidden directory was skipped
	hiddenIndexPath := filepath.Join(hiddenDir, ".hidden.md")
	if _, err := os.Stat(hiddenIndexPath); err == nil {
		t.Error("Hidden directory should not have an index file created")
	}

	// Verify that normal directory has index file
	normalIndexPath := filepath.Join(normalDir, "normal.md")
	if _, err := os.Stat(normalIndexPath); os.IsNotExist(err) {
		t.Error("Normal directory should have an index file created")
	}
}

func TestIndexator_Start_WithExistingIndexFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create directory with existing index file
	testDir := filepath.Join(tempDir, "testdir")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create existing index file
	existingIndex := filepath.Join(testDir, "testdir.md")
	err = os.WriteFile(existingIndex, []byte("# Existing index\n[[old.md]]"), 0644)
	if err != nil {
		t.Fatalf("Failed to create existing index file: %v", err)
	}

	// Create some files in the directory
	testFile := filepath.Join(testDir, "newfile.md")
	err = os.WriteFile(testFile, []byte("# New file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	indexator := NewIndexator(tempDir)
	err = indexator.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// The existing index file should be skipped (not overwritten)
	content, err := os.ReadFile(existingIndex)
	if err != nil {
		t.Fatalf("Failed to read existing index file: %v", err)
	}

	// Should still contain the original content
	if !strings.Contains(string(content), "[[old.md]]") {
		t.Error("Existing index file should not have been modified")
	}
}

func TestIndexator_createIndexFile(t *testing.T) {
	tempDir := t.TempDir()

	indexator := &Indexator{vaultPath: tempDir}

	links := []string{
		"[[file1.md]]",
		"[[file2.md]]",
		"[[subdir/subfile.md]]",
	}

	// Test creating index file in subdirectory
	subDir := filepath.Join(tempDir, "testdir")
	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	err = indexator.createIndexFile(subDir, links)
	if err != nil {
		t.Fatalf("createIndexFile() failed: %v", err)
	}

	// Verify file was created
	indexPath := filepath.Join(subDir, "testdir.md")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("Index file was not created")
	}

	// Verify content
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index file: %v", err)
	}

	contentStr := string(content)
	for _, link := range links {
		if !strings.Contains(contentStr, link) {
			t.Errorf("Index file should contain link %s", link)
		}
	}

	// Test creating index file in root directory
	rootLinks := []string{"[[rootfile.md]]"}
	err = indexator.createIndexFile(tempDir, rootLinks)
	if err != nil {
		t.Fatalf("createIndexFile() failed for root: %v", err)
	}

	// Verify root index file was created
	rootIndexPath := filepath.Join(tempDir, "index.md")
	if _, err := os.Stat(rootIndexPath); os.IsNotExist(err) {
		t.Error("Root index file was not created")
	}
}

func TestIndexator_collectDirectories(t *testing.T) {
	tempDir := t.TempDir()

	// Create test directory structure
	testDirs := []string{
		"level1",
		"level1/level2",
		"level1/level2/level3",
		"another",
	}

	for _, dir := range testDirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	indexator := &Indexator{vaultPath: tempDir}
	directories, err := indexator.CollectDirectories()
	if err != nil {
		t.Fatalf("collectDirectories() failed: %v", err)
	}

	// Should include all directories including root
	expectedCount := len(testDirs) + 1 // +1 for root directory
	if len(directories) != expectedCount {
		t.Errorf("Expected %d directories, got %d", expectedCount, len(directories))
	}

	// Verify all expected directories are present
	dirMap := make(map[string]bool)
	for _, dir := range directories {
		dirMap[dir] = true
	}

	for _, expectedDir := range testDirs {
		if !dirMap[expectedDir] {
			t.Errorf("Expected directory %s not found in collected directories", expectedDir)
		}
	}

	// Root directory should be present
	if !dirMap["."] {
		t.Error("Root directory should be present in collected directories")
	}
}

func TestIndexator_Start_WithNestedStructure(t *testing.T) {
	tempDir := t.TempDir()

	// Create a complex nested structure
	structure := map[string]string{
		"notes/2023/january/meeting1.md": "# Meeting 1",
		"notes/2023/january/meeting2.md": "# Meeting 2",
		"notes/2023/february/report.md":  "# Report",
		"notes/2024/planning/goals.md":   "# Goals",
		"docs/readme.md":                 "# Readme",
		"docs/api/reference.md":          "# API Reference",
	}

	for filePath, content := range structure {
		fullPath := filepath.Join(tempDir, filePath)
		dir := filepath.Dir(fullPath)

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", filePath, err)
		}
	}

	indexator := NewIndexator(tempDir)
	err := indexator.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Verify that index files were created at all levels
	expectedIndexFiles := []string{
		"notes/notes.md",
		"notes/2023/2023.md",
		"notes/2023/january/january.md",
		"notes/2023/february/february.md",
		"notes/2024/2024.md",
		"notes/2024/planning/planning.md",
		"docs/docs.md",
		"docs/api/api.md",
	}

	for _, indexFile := range expectedIndexFiles {
		fullPath := filepath.Join(tempDir, indexFile)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected index file %s was not created", indexFile)
		}
	}

	// Verify that the deepest index files contain the correct links
	januaryIndexPath := filepath.Join(tempDir, "notes/2023/january/january.md")
	content, err := os.ReadFile(januaryIndexPath)
	if err != nil {
		t.Fatalf("Failed to read january index file: %v", err)
	}

	contentStr := string(content)
	expectedLinks := []string{
		"[[notes/2023/january/meeting1.md]]",
		"[[notes/2023/january/meeting2.md]]",
	}

	for _, link := range expectedLinks {
		if !strings.Contains(contentStr, link) {
			t.Errorf("January index should contain link %s", link)
		}
	}
}
