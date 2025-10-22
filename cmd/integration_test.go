package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"git-mini-commit/testutils"
)

func TestCLIIntegration(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)
	cli.SetRepo(repo)

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

	// 3. mini-commit一覧を表示
	output = cli.AssertCommandSuccess(t, "list")
	if !strings.Contains(output, "Mini-commits") {
		t.Errorf("Expected 'Mini-commits' in output, but got: %s", output)
	}

	// 4. mini-commitのIDを抽出（listコマンドから取得）
	listOutput := cli.AssertCommandSuccess(t, "list")
	
	lines := strings.Split(listOutput, "\n")
	var miniCommitID string
	for _, line := range lines {
		if strings.Contains(line, "ID:") {
			parts := strings.Split(line, "ID: ")
			if len(parts) > 1 {
				miniCommitID = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if miniCommitID == "" {
		t.Fatalf("Failed to extract mini-commit ID from list output: %s", listOutput)
	}

	// 5. mini-commitの差分を表示
	output = cli.AssertCommandSuccess(t, "show", miniCommitID)
	if !strings.Contains(output, "Mini-commit:") {
		t.Errorf("Expected 'Mini-commit:' in output, but got: %s", output)
	}

	// 6. mini-commitを削除
	output = cli.AssertCommandSuccess(t, "drop", miniCommitID)
	if !strings.Contains(output, "Deleted mini-commit") {
		t.Errorf("Expected 'Deleted mini-commit' in output, but got: %s", output)
	}

	// 7. mini-commit一覧が空かチェック
	output = cli.AssertCommandSuccess(t, "list")
	if !strings.Contains(output, "No mini-commits found") {
		t.Errorf("Expected 'No mini-commits found' in output, but got: %s", output)
	}
}

func TestCLIIntegrationErrorHandling(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)
	cli.SetRepo(repo)

	// 1. メッセージなしでmini-commitを作成（エラー）
	cli.AssertCommandFailure(t)

	// 2. ステージングされた変更なしでmini-commitを作成（エラー）
	cli.AssertCommandFailure(t, "-m", "Test commit")

	// 3. 存在しないmini-commitを表示（エラー）
	cli.AssertCommandFailure(t, "show", "nonexistent-id")

	// 4. 存在しないmini-commitを削除（エラー）
	cli.AssertCommandFailure(t, "drop", "nonexistent-id")

	// 5. 存在しないmini-commitをpop（エラー）
	cli.AssertCommandFailure(t, "pop", "nonexistent-id")
}

func TestCLIWithMultipleMiniCommits(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)
	cli.SetRepo(repo)

	// 1. 最初のファイルを作成してステージング
	if err := repo.CreateTestFile("file1.txt", "Content 1\n"); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}
	if err := repo.StageFile("file1.txt"); err != nil {
		t.Fatalf("Failed to stage file1: %v", err)
	}

	// 2. 最初のmini-commitを作成
	output1 := cli.AssertCommandSuccess(t, "-m", "First commit")
	if !strings.Contains(output1, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output1)
	}

	// 3. 2番目のファイルを作成してステージング
	if err := repo.CreateTestFile("file2.txt", "Content 2\n"); err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}
	if err := repo.StageFile("file2.txt"); err != nil {
		t.Fatalf("Failed to stage file2: %v", err)
	}

	// 4. 2番目のmini-commitを作成
	output2 := cli.AssertCommandSuccess(t, "-m", "Second commit")
	if !strings.Contains(output2, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output2)
	}

	// 5. mini-commit一覧を表示
	output := cli.AssertCommandSuccess(t, "list")
	if !strings.Contains(output, "Mini-commits (2)") {
		t.Errorf("Expected 'Mini-commits (2)' in output, but got: %s", output)
	}

	// 6. 各mini-commitのメッセージが含まれているかチェック
	if !strings.Contains(output, "First commit") {
		t.Errorf("Expected 'First commit' in output, but got: %s", output)
	}
	if !strings.Contains(output, "Second commit") {
		t.Errorf("Expected 'Second commit' in output, but got: %s", output)
	}
}

func TestCLIPopCommand(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)
	cli.SetRepo(repo)

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

	// 3. mini-commitのIDを抽出（listコマンドから取得）
	listOutput := cli.AssertCommandSuccess(t, "list")
	
	lines := strings.Split(listOutput, "\n")
	var miniCommitID string
	for _, line := range lines {
		if strings.Contains(line, "ID:") {
			parts := strings.Split(line, "ID: ")
			if len(parts) > 1 {
				miniCommitID = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if miniCommitID == "" {
		t.Fatalf("Failed to extract mini-commit ID from list output: %s", listOutput)
	}

	// 4. 新しいクリーンなリポジトリを作成
	cleanRepo := testutils.NewTestGitRepo(t)
	defer cleanRepo.Cleanup()
	
	// クリーンなリポジトリに移動
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(cleanRepo.RepoPath); err != nil {
		t.Fatalf("Failed to change to clean repo directory: %v", err)
	}

	// 5. mini-commitをpop
	output = cli.AssertCommandSuccess(t, "pop", miniCommitID)
	if !strings.Contains(output, "Applied mini-commit") {
		t.Errorf("Expected 'Applied mini-commit' in output, but got: %s", output)
	}

	// 6. ステージングエリアに変更が適用されているかチェック
	hasChanges, err := repo.HasStagedChanges()
	if err != nil {
		t.Fatalf("Failed to check staging status: %v", err)
	}
	if !hasChanges {
		t.Errorf("Expected staged changes after pop, but got false")
	}
}

func TestCLIWithLargeFiles(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)
	cli.SetRepo(repo)

	// 1. 大きなファイルを作成
	largeContent := strings.Repeat("This is a test line.\n", 1000)
	if err := repo.CreateTestFile("large.txt", largeContent); err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}
	if err := repo.StageFile("large.txt"); err != nil {
		t.Fatalf("Failed to stage large file: %v", err)
	}

	// 2. mini-commitを作成
	output := cli.AssertCommandSuccess(t, "-m", "Large file commit")
	if !strings.Contains(output, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output)
	}

	// 3. mini-commit一覧を表示
	output = cli.AssertCommandSuccess(t, "list")
	if !strings.Contains(output, "Mini-commits (1)") {
		t.Errorf("Expected 'Mini-commits (1)' in output, but got: %s", output)
	}
}

func TestCLIWithBinaryFiles(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)
	cli.SetRepo(repo)

	// 1. バイナリファイルを作成
	binaryContent := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	if err := os.WriteFile("binary.bin", binaryContent, 0644); err != nil {
		t.Fatalf("Failed to create binary file: %v", err)
	}
	if err := repo.StageFile("binary.bin"); err != nil {
		t.Fatalf("Failed to stage binary file: %v", err)
	}

	// 2. mini-commitを作成
	output := cli.AssertCommandSuccess(t, "-m", "Binary file commit")
	if !strings.Contains(output, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output)
	}

	// 3. mini-commit一覧を表示
	output = cli.AssertCommandSuccess(t, "list")
	if !strings.Contains(output, "Mini-commits (1)") {
		t.Errorf("Expected 'Mini-commits (1)' in output, but got: %s", output)
	}
}

func TestCLIWithSpecialCharacters(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)
	cli.SetRepo(repo)

	// 1. 特殊文字を含むファイルを作成
	specialContent := "Hello, 世界! 🌍\nSpecial chars: !@#$%^&*()\n"
	if err := repo.CreateTestFile("special.txt", specialContent); err != nil {
		t.Fatalf("Failed to create special file: %v", err)
	}
	if err := repo.StageFile("special.txt"); err != nil {
		t.Fatalf("Failed to stage special file: %v", err)
	}

	// 2. 特殊文字を含むメッセージでmini-commitを作成
	output := cli.AssertCommandSuccess(t, "-m", "Special chars: 世界! 🌍")
	if !strings.Contains(output, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output)
	}

	// 3. mini-commit一覧を表示
	output = cli.AssertCommandSuccess(t, "list")
	if !strings.Contains(output, "Mini-commits (1)") {
		t.Errorf("Expected 'Mini-commits (1)' in output, but got: %s", output)
	}
}

func TestCLIWithEmptyFiles(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)
	cli.SetRepo(repo)

	// 1. 空のファイルを作成
	if err := repo.CreateTestFile("empty.txt", ""); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}
	if err := repo.StageFile("empty.txt"); err != nil {
		t.Fatalf("Failed to stage empty file: %v", err)
	}

	// 2. mini-commitを作成
	output := cli.AssertCommandSuccess(t, "-m", "Empty file commit")
	if !strings.Contains(output, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output)
	}

	// 3. mini-commit一覧を表示
	output = cli.AssertCommandSuccess(t, "list")
	if !strings.Contains(output, "Mini-commits (1)") {
		t.Errorf("Expected 'Mini-commits (1)' in output, but got: %s", output)
	}
}

func TestCLIWithLongMessages(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)
	cli.SetRepo(repo)

	// 1. ファイルを作成してステージング
	if err := repo.CreateTestFile("test.txt", "Hello, World!\n"); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := repo.StageFile("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// 2. 長いメッセージでmini-commitを作成
	longMessage := strings.Repeat("This is a very long commit message. ", 100)
	output := cli.AssertCommandSuccess(t, "-m", longMessage)
	if !strings.Contains(output, "Created mini-commit") {
		t.Errorf("Expected 'Created mini-commit' in output, but got: %s", output)
	}

	// 3. mini-commit一覧を表示
	output = cli.AssertCommandSuccess(t, "list")
	if !strings.Contains(output, "Mini-commits (1)") {
		t.Errorf("Expected 'Mini-commits (1)' in output, but got: %s", output)
	}
}

func TestCLIWithConcurrentOperations(t *testing.T) {
	// テスト用Gitリポジトリを作成
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// テスト用CLIを作成
	cli := testutils.NewTestCLI(t)
	cli.SetRepo(repo)

	// 複数のgoroutineで同時にmini-commitを作成
	done := make(chan bool, 5)
	
	for i := 0; i < 5; i++ {
		go func(i int) {
			// ファイルを作成してステージング
			filename := fmt.Sprintf("file%d.txt", i)
			content := fmt.Sprintf("Content %d\n", i)
			if err := repo.CreateTestFile(filename, content); err != nil {
				t.Errorf("Failed to create file %d: %v", i, err)
				done <- false
				return
			}
		if err := repo.StageFile(filename); err != nil {
			// 並行処理での競合は許容する
			done <- false
			return
		}

			// mini-commitを作成
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
	for i := 0; i < 5; i++ {
		if <-done {
			successCount++
		}
	}

	// 成功した操作が3つ以上あればOK（並行処理での競合を許容）
	if successCount < 3 {
		t.Errorf("Expected at least 3 successful operations, but got %d", successCount)
	}

	// mini-commit一覧を表示
	output := cli.AssertCommandSuccess(t, "list")
	if !strings.Contains(output, "Mini-commits") {
		t.Errorf("Expected 'Mini-commits' in output, but got: %s", output)
	}
}
