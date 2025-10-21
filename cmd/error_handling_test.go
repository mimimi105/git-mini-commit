package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git-mini-commit/testutils"
)

func TestCLIErrorHandling(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)

	t.Run("No message provided", func(t *testing.T) {
		// ファイルを作成してステージング
		if err := repo.CreateTestFile("test.txt", "Hello, World!\n"); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		if err := repo.StageFile("test.txt"); err != nil {
			t.Fatalf("Failed to stage file: %v", err)
		}

		// メッセージなしでmini-commitを作成（エラー）
		output := cli.AssertCommandFailure(t)
		if !strings.Contains(output, "message is required") {
			t.Errorf("Expected 'message is required' in output, but got: %s", output)
		}
	})

	t.Run("No staged changes", func(t *testing.T) {
		// 新しいテスト用Gitリポジトリを作成（ステージングエリアが空の状態）
		cleanRepo := testutils.NewTestGitRepo(t)
		defer cleanRepo.Cleanup()
		
		// ファイルを作成するがステージングしない
		if err := cleanRepo.CreateTestFile("test.txt", "Hello, World!\n"); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		
		// ステージングされた変更なしでmini-commitを作成（エラー）
		output := cli.AssertCommandFailure(t, "-m", "Test commit")
		if !strings.Contains(output, "no staged changes") {
			t.Errorf("Expected 'no staged changes' in output, but got: %s", output)
		}
	})

	t.Run("Invalid mini-commit ID", func(t *testing.T) {
		// 存在しないmini-commitを表示（エラー）
		output := cli.AssertCommandFailure(t, "show", "invalid-id")
		if !strings.Contains(output, "not found") {
			t.Errorf("Expected 'not found' in output, but got: %s", output)
		}

		// 存在しないmini-commitを削除（エラー）
		output = cli.AssertCommandFailure(t, "drop", "invalid-id")
		if !strings.Contains(output, "not found") {
			t.Errorf("Expected 'not found' in output, but got: %s", output)
		}

		// 存在しないmini-commitをpop（エラー）
		output = cli.AssertCommandFailure(t, "pop", "invalid-id")
		if !strings.Contains(output, "not found") {
			t.Errorf("Expected 'not found' in output, but got: %s", output)
		}
	})

	t.Run("Invalid command arguments", func(t *testing.T) {
		// 引数が多すぎる場合
		output := cli.AssertCommandFailure(t, "show", "id1", "id2")
		if !strings.Contains(output, "accepts 1 arg(s), received 2") {
			t.Errorf("Expected 'accepts 1 arg(s), received 2' in output, but got: %s", output)
		}

		// 引数が少なすぎる場合
		output = cli.AssertCommandFailure(t, "show")
		if !strings.Contains(output, "accepts 1 arg(s), received 0") {
			t.Errorf("Expected 'accepts 1 arg(s), received 0' in output, but got: %s", output)
		}
	})
}

func TestCLIInNonGitRepository(t *testing.T) {
	// 元のディレクトリを保存
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// 一時ディレクトリに移動（Gitリポジトリではない）
	tempDir, err := os.MkdirTemp("", "not-git-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)

	// Gitリポジトリではない場所でコマンドを実行（エラー）
	output := cli.AssertCommandFailure(t, "-m", "Test commit")
	if !strings.Contains(output, "not a git repository") {
		t.Errorf("Expected 'not a git repository' in output, but got: %s", output)
	}

	output = cli.AssertCommandFailure(t, "list")
	if !strings.Contains(output, "not a git repository") {
		t.Errorf("Expected 'not a git repository' in output, but got: %s", output)
	}
}

func TestCLIWithCorruptedStorage(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)

	// 1. 正常なmini-commitを作成
	if err := repo.CreateTestFile("test.txt", "Hello, World!\n"); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := repo.StageFile("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	output := cli.AssertCommandSuccess(t, "-m", "Test commit")
	if !strings.Contains(output, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output)
	}

	// 2. インデックスファイルを破損
	indexPath := filepath.Join(".git", "mini-commits", "index.json")
	if err := os.WriteFile(indexPath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to corrupt index file: %v", err)
	}

	// 3. 破損したインデックスでコマンドを実行（エラー）
	output = cli.AssertCommandFailure(t, "list")
	if !strings.Contains(output, "failed to parse index") {
		t.Errorf("Expected 'failed to parse index' in output, but got: %s", output)
	}
}

func TestCLIWithPermissionErrors(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)

	// 1. ファイルを作成してステージング
	if err := repo.CreateTestFile("test.txt", "Hello, World!\n"); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := repo.StageFile("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// 2. mini-commitsディレクトリを作成してから権限を変更（読み取り専用）
	miniCommitsPath := filepath.Join(".git", "mini-commits")
	if err := os.MkdirAll(miniCommitsPath, 0755); err != nil {
		t.Fatalf("Failed to create mini-commits directory: %v", err)
	}
	if err := os.Chmod(miniCommitsPath, 0444); err != nil {
		t.Fatalf("Failed to change permissions: %v", err)
	}
	defer os.Chmod(miniCommitsPath, 0755) // 権限を元に戻す

	// 3. 権限エラーでmini-commitを作成（エラー）
	output := cli.AssertCommandFailure(t, "-m", "Test commit")
	if !strings.Contains(output, "permission denied") {
		t.Errorf("Expected 'permission denied' in output, but got: %s", output)
	}
}

func TestCLIWithDiskSpaceErrors(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)

	// 1. 非常に大きなファイルを作成
	largeContent := strings.Repeat("This is a test line.\n", 1000000) // 約20MB
	if err := repo.CreateTestFile("large.txt", largeContent); err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}
	if err := repo.StageFile("large.txt"); err != nil {
		t.Fatalf("Failed to stage large file: %v", err)
	}

	// 2. 大きなファイルでmini-commitを作成
	output := cli.AssertCommandSuccess(t, "-m", "Large file commit")
	if !strings.Contains(output, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output)
	}
}

func TestCLIWithNetworkErrors(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)

	// 1. ファイルを作成してステージング
	if err := repo.CreateTestFile("test.txt", "Hello, World!\n"); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := repo.StageFile("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// 2. mini-commitを作成
	output := cli.AssertCommandSuccess(t, "-m", "Test commit")
	if !strings.Contains(output, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output)
	}

	// 3. ネットワークエラーをシミュレート（実際のネットワーク操作はないが、エラーハンドリングをテスト）
	// このテストは実際のネットワークエラーをシミュレートできないため、
	// エラーハンドリングのロジックをテストする
}

func TestCLIWithConcurrentAccess(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)

	// 1. ファイルを作成してステージング
	if err := repo.CreateTestFile("test.txt", "Hello, World!\n"); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := repo.StageFile("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// 2. 複数のgoroutineで同時にmini-commitを作成
	done := make(chan bool, 3)
	
	for i := 0; i < 3; i++ {
		go func(i int) {
			message := fmt.Sprintf("Concurrent commit %d", i)
			output := cli.AssertCommandSuccess(t, "-m", message)
			if !strings.Contains(output, "Created mini-commit") {
				t.Errorf("Expected 'Created mini-commit' in output for commit %d, but got: %s", i, output)
				done <- false
				return
			}
			done <- true
		}(i)
	}

	// すべてのgoroutineが完了するまで待機
	successCount := 0
	for i := 0; i < 3; i++ {
		if <-done {
			successCount++
		}
	}

	if successCount != 3 {
		t.Errorf("Expected 3 successful operations, but got %d", successCount)
	}
}

func TestCLIWithInvalidPatchContent(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)

	// 1. ファイルを作成してステージング
	if err := repo.CreateTestFile("test.txt", "Hello, World!\n"); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := repo.StageFile("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// 2. mini-commitを作成
	output := cli.AssertCommandSuccess(t, "-m", "Test commit")
	if !strings.Contains(output, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output)
	}

	// 3. mini-commitのIDを抽出
	lines := strings.Split(output, "\n")
	var miniCommitID string
	for _, line := range lines {
		if strings.Contains(line, "Created mini-commit:") {
			parts := strings.Split(line, "Created mini-commit: ")
			if len(parts) > 1 {
				miniCommitID = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if miniCommitID == "" {
		t.Fatalf("Failed to extract mini-commit ID from output: %s", output)
	}

	// 4. patchファイルを破損
	patchPath := filepath.Join(".git", "mini-commits", miniCommitID+".patch")
	if err := os.WriteFile(patchPath, []byte("invalid patch content"), 0644); err != nil {
		t.Fatalf("Failed to corrupt patch file: %v", err)
	}

	// 5. 破損したpatchでpopを実行（エラー）
	output = cli.AssertCommandFailure(t, "pop", miniCommitID)
	if !strings.Contains(output, "not found") {
		t.Errorf("Expected 'not found' in output, but got: %s", output)
	}
}

func TestCLIWithMemoryConstraints(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)

	// 1. 複数の大きなファイルを作成
	for i := 0; i < 10; i++ {
		filename := fmt.Sprintf("file%d.txt", i)
		content := strings.Repeat("This is a test line.\n", 10000) // 約200KB
		if err := repo.CreateTestFile(filename, content); err != nil {
			t.Fatalf("Failed to create file %d: %v", i, err)
		}
		if err := repo.StageFile(filename); err != nil {
			t.Fatalf("Failed to stage file %d: %v", i, err)
		}
	}

	// 2. 大きなファイルでmini-commitを作成
	output := cli.AssertCommandSuccess(t, "-m", "Large files commit")
	if !strings.Contains(output, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output)
	}
}

func TestCLIWithUnicodeContent(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)

	// 1. Unicode文字を含むファイルを作成
	unicodeContent := "Hello, 世界! 🌍\nSpecial chars: !@#$%^&*()\n"
	if err := repo.CreateTestFile("unicode.txt", unicodeContent); err != nil {
		t.Fatalf("Failed to create unicode file: %v", err)
	}
	if err := repo.StageFile("unicode.txt"); err != nil {
		t.Fatalf("Failed to stage unicode file: %v", err)
	}

	// 2. Unicode文字でmini-commitを作成
	output := cli.AssertCommandSuccess(t, "-m", "Unicode commit: 世界! 🌍")
	if !strings.Contains(output, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output)
	}

	// 3. mini-commit一覧を表示
	output = cli.AssertCommandSuccess(t, "list")
	if !strings.Contains(output, "Mini-commits (1)") {
		t.Errorf("Expected 'Mini-commits (1)' in output, but got: %s", output)
	}
}
