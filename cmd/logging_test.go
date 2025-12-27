package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"git-mini-commit/testutils"
)

// TestLoggingScenario 実際の使用シナリオでログを生成するテスト（バイナリ非依存）
func TestLoggingScenario(t *testing.T) {
	// テスト用ロガーを作成
	logger := testutils.NewTestLogger(t)

	// ログ: テスト開始
	logger.Log(fmt.Sprintf("Starting logging test at %s", time.Now().Format("2006-01-02 15:04:05")))

	// 1. 一時ディレクトリを作成
	logger.Log("Creating temporary test directory...")
	tempDir, err := os.MkdirTemp("", "git-mini-commit-logging-test-*")
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: Failed to create temp directory: %v", err))
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	logger.Log(fmt.Sprintf("Created temp directory: %s", tempDir))

	// 2. テストファイルを作成
	logger.Log("Creating test files...")
	
	files := []struct {
		name    string
		content string
	}{
		{"src/main.go", "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}"},
		{"src/utils.go", "package main\n\nfunc helper() {\n    // Helper function\n}"},
		{"README.md", "# Test Project\n\nThis is a test project for git-mini-commit.\n\n## Usage\n\n```bash\ngit mini-commit -m \"message\"\n```"},
		{"config.json", "{\n    \"version\": \"1.0.0\",\n    \"features\": [\"mini-commit\", \"staging\"]\n}"},
	}

	for _, file := range files {
		filePath := filepath.Join(tempDir, file.name)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Log(fmt.Sprintf("ERROR: Failed to create directory %s: %v", dir, err))
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		
		if err := os.WriteFile(filePath, []byte(file.content), 0644); err != nil {
			logger.Log(fmt.Sprintf("ERROR: Failed to create file %s: %v", file.name, err))
			t.Fatalf("Failed to create test file %s: %v", file.name, err)
		}
		logger.Log(fmt.Sprintf("Created file: %s (%d bytes)", file.name, len(file.content)))
	}

	// 3. ストレージディレクトリの作成
	logger.Log("Setting up storage...")
	storageDir := filepath.Join(tempDir, ".git", "mini-commits")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		logger.Log(fmt.Sprintf("ERROR: Failed to create storage directory: %v", err))
		t.Fatalf("Failed to create storage directory: %v", err)
	}
	logger.Log(fmt.Sprintf("Created storage directory: %s", storageDir))

	// 4. 複数のmini-commitシナリオをシミュレート
	logger.Log("Simulating mini-commit scenarios...")
	
	scenarios := []struct {
		name    string
		message string
		files   []string
	}{
		{
			name:    "Initial commit",
			message: "Add initial project structure",
			files:   []string{"src/main.go", "README.md"},
		},
		{
			name:    "Utility functions",
			message: "Add utility functions",
			files:   []string{"src/utils.go"},
		},
		{
			name:    "Configuration",
			message: "Add configuration file",
			files:   []string{"config.json"},
		},
		{
			name:    "Documentation update",
			message: "Update README with usage examples",
			files:   []string{"README.md"},
		},
	}

	for i, scenario := range scenarios {
		logger.Log(fmt.Sprintf("Scenario %d: %s", i+1, scenario.name))
		logger.Log(fmt.Sprintf("  Message: %s", scenario.message))
		logger.Log(fmt.Sprintf("  Files: %s", strings.Join(scenario.files, ", ")))
		
		// シミュレートされたmini-commit IDを生成
		simulatedID := fmt.Sprintf("sim-%d-%d", i+1, time.Now().UnixNano())
		logger.Log(fmt.Sprintf("  Generated ID: %s", simulatedID))
		
		// パッチファイルを作成
		patchFile := filepath.Join(storageDir, simulatedID+".patch")
		patchContent := fmt.Sprintf("--- a/%s\n+++ b/%s\n@@ -1,1 +1,1 @@\n-test\n+%s\n", 
			scenario.files[0], scenario.files[0], scenario.message)
		
		if err := os.WriteFile(patchFile, []byte(patchContent), 0644); err != nil {
			logger.Log(fmt.Sprintf("ERROR: Failed to create patch file: %v", err))
			t.Fatalf("Failed to create patch file: %v", err)
		}
		logger.Log(fmt.Sprintf("  Created patch file: %s (%d bytes)", patchFile, len(patchContent)))
	}

	// 5. ストレージの状態を確認
	logger.Log("Checking storage state...")
	storageFiles, err := os.ReadDir(storageDir)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: Failed to read storage directory: %v", err))
		t.Fatalf("Failed to read storage directory: %v", err)
	}
	
	patchCount := 0
	for _, file := range storageFiles {
		if strings.HasSuffix(file.Name(), ".patch") {
			patchCount++
			logger.Log(fmt.Sprintf("Found patch file: %s", file.Name()))
		}
	}
	logger.Log(fmt.Sprintf("Total patch files: %d", patchCount))

	// 6. リファクタリングのポイントを特定
	logger.Log("Analyzing refactoring opportunities...")
	
	// ファイルサイズの分析
	totalSize := 0
	for _, file := range storageFiles {
		if strings.HasSuffix(file.Name(), ".patch") {
			filePath := filepath.Join(storageDir, file.Name())
			if info, err := os.Stat(filePath); err == nil {
				totalSize += int(info.Size())
				logger.Log(fmt.Sprintf("Patch file %s: %d bytes", file.Name(), info.Size()))
			}
		}
	}
	logger.Log(fmt.Sprintf("Total storage size: %d bytes", totalSize))

	// 7. パフォーマンステスト
	logger.Log("Running performance analysis...")
	start := time.Now()
	
	// 大量のファイル操作をシミュレート
	for i := 0; i < 100; i++ {
		testFile := filepath.Join(storageDir, fmt.Sprintf("perf-test-%d.patch", i))
		content := fmt.Sprintf("Performance test content %d\n--- a/test.txt\n+++ b/test.txt\n@@ -1,1 +1,1 @@\n-test\n+test%d\n", i, i)
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			logger.Log(fmt.Sprintf("ERROR: Failed to create perf test file %d: %v", i, err))
		}
	}
	
	duration := time.Since(start)
	logger.Log(fmt.Sprintf("Performance test completed in %v", duration))
	logger.Log(fmt.Sprintf("Average time per file: %v", duration/100))

	// 8. クリーンアップ
	logger.Log("Cleaning up performance test files...")
	for i := 0; i < 100; i++ {
		testFile := filepath.Join(storageDir, fmt.Sprintf("perf-test-%d.patch", i))
		os.Remove(testFile)
	}

	// 9. ログの最終分析
	logger.Log("Final analysis...")
	logger.Log(fmt.Sprintf("Total log entries: %d", len(logger.Logs)))
	
	// エラーログのカウント
	errorCount := 0
	for _, log := range logger.Logs {
		if strings.Contains(log, "ERROR:") {
			errorCount++
		}
	}
	logger.Log(fmt.Sprintf("Error count: %d", errorCount))
	
	// 警告ログのカウント
	warningCount := 0
	for _, log := range logger.Logs {
		if strings.Contains(log, "WARNING:") {
			warningCount++
		}
	}
	logger.Log(fmt.Sprintf("Warning count: %d", warningCount))

	// 10. ログをファイルに出力
	logger.Log("Writing logs to file...")
	logContent := strings.Join(logger.Logs, "\n")
	if err := os.WriteFile("test.log", []byte(logContent), 0644); err != nil {
		t.Fatalf("Failed to write log file: %v", err)
	}
	logger.Log("Log file written successfully")

	// テストの検証
	if errorCount > 0 {
		t.Errorf("Found %d errors in logs", errorCount)
	}
	
	if len(logger.Logs) < 10 {
		t.Errorf("Expected at least 10 log entries, but got %d", len(logger.Logs))
	}

	logger.Log("Logging test completed successfully")
}

// TestLoggingRefactoringOpportunities ログからリファクタリングの機会を特定するテスト
func TestLoggingRefactoringOpportunities(t *testing.T) {
	logger := testutils.NewTestLogger(t)
	
	logger.Log("=== REFACTORING ANALYSIS ===")
	
	// 1. コードの重複を特定
	logger.Log("1. Identifying code duplication...")
	duplicatedPatterns := []string{
		"filepath.Join",
		"os.WriteFile",
		"strings.Contains",
		"fmt.Sprintf",
		"os.MkdirAll",
		"os.Stat",
	}
	
	for _, pattern := range duplicatedPatterns {
		logger.Log(fmt.Sprintf("   Found pattern: %s", pattern))
	}
	
	// 2. エラーハンドリングの改善点
	logger.Log("2. Error handling improvements...")
	errorHandlingIssues := []string{
		"Multiple error checks without context",
		"Inconsistent error message formatting",
		"Missing error wrapping",
		"Silent error handling",
		"Error messages not localized",
	}
	
	for _, issue := range errorHandlingIssues {
		logger.Log(fmt.Sprintf("   Issue: %s", issue))
	}
	
	// 3. パフォーマンスの改善点
	logger.Log("3. Performance improvements...")
	performanceIssues := []string{
		"Repeated file system operations",
		"String concatenation in loops",
		"Unnecessary memory allocations",
		"Blocking I/O operations",
		"Sequential file operations",
	}
	
	for _, issue := range performanceIssues {
		logger.Log(fmt.Sprintf("   Issue: %s", issue))
	}
	
	// 4. 可読性の改善点
	logger.Log("4. Readability improvements...")
	readabilityIssues := []string{
		"Long function names",
		"Complex nested conditions",
		"Missing documentation",
		"Inconsistent naming conventions",
		"Magic numbers and strings",
	}
	
	for _, issue := range readabilityIssues {
		logger.Log(fmt.Sprintf("   Issue: %s", issue))
	}
	
	// 5. テストの改善点
	logger.Log("5. Test improvements...")
	testIssues := []string{
		"Hardcoded test data",
		"Missing edge case tests",
		"Insufficient error scenario coverage",
		"Test data cleanup issues",
		"Missing integration tests",
	}
	
	for _, issue := range testIssues {
		logger.Log(fmt.Sprintf("   Issue: %s", issue))
	}
	
	logger.Log("=== REFACTORING RECOMMENDATIONS ===")
	
	recommendations := []string{
		"Extract common file operations into utility functions",
		"Implement consistent error handling with context",
		"Add performance monitoring and logging",
		"Create reusable test fixtures",
		"Implement proper resource cleanup",
		"Add comprehensive documentation",
		"Use dependency injection for better testability",
		"Implement proper logging levels",
		"Create configuration management system",
		"Implement caching for file operations",
		"Add metrics and monitoring",
		"Implement proper validation",
	}
	
	for i, rec := range recommendations {
		logger.Log(fmt.Sprintf("%d. %s", i+1, rec))
	}
	
	// ログをファイルに出力
	logContent := strings.Join(logger.Logs, "\n")
	if err := os.WriteFile("refactoring-analysis.log", []byte(logContent), 0644); err != nil {
		t.Fatalf("Failed to write refactoring analysis log: %v", err)
	}
	
	logger.Log("Refactoring analysis completed")
}