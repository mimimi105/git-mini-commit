package testutils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestGitRepo テスト用のGitリポジトリ管理
type TestGitRepo struct {
	RepoPath string
	OriginalDir string
}

// NewTestGitRepo テスト用のGitリポジトリを作成
func NewTestGitRepo(t *testing.T) *TestGitRepo {
	// 一時ディレクトリを作成
	repoPath, err := os.MkdirTemp("", "git-mini-commit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// 元のディレクトリを保存
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// テスト用ディレクトリに移動
	if err := os.Chdir(repoPath); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Gitリポジトリを初期化
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// ユーザー設定（テスト用）
	exec.Command("git", "config", "user.name", "Test User").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()

	return &TestGitRepo{
		RepoPath: repoPath,
		OriginalDir: originalDir,
	}
}

// Cleanup テスト用リポジトリをクリーンアップ
func (r *TestGitRepo) Cleanup() {
	// 元のディレクトリに戻る
	os.Chdir(r.OriginalDir)
	
	// 一時ディレクトリを削除
	os.RemoveAll(r.RepoPath)
}

// CreateTestFile テスト用ファイルを作成
func (r *TestGitRepo) CreateTestFile(filename, content string) error {
	filePath := filepath.Join(r.RepoPath, filename)
	return os.WriteFile(filePath, []byte(content), 0644)
}

// StageFile ファイルをステージング
func (r *TestGitRepo) StageFile(filename string) error {
	return exec.Command("git", "add", filename).Run()
}

// CommitFile ファイルをコミット
func (r *TestGitRepo) CommitFile(message string) error {
	return exec.Command("git", "commit", "-m", message).Run()
}

// GetStagedChanges ステージングされた変更を取得
func (r *TestGitRepo) GetStagedChanges() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	return stdout.String(), err
}

// HasStagedChanges ステージングされた変更があるかチェック
func (r *TestGitRepo) HasStagedChanges() (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

// TestStorage テスト用ストレージ管理（簡易版）
type TestStorage struct {
	BasePath string
}

// NewTestStorage テスト用ストレージを作成
func NewTestStorage(t *testing.T) *TestStorage {
	// テスト用のmini-commitsディレクトリを作成
	basePath := filepath.Join(".git", "mini-commits")
	if err := os.MkdirAll(basePath, 0755); err != nil {
		t.Fatalf("Failed to create mini-commits directory: %v", err)
	}

	return &TestStorage{
		BasePath: basePath,
	}
}

// CreateTestMiniCommit テスト用mini-commitを作成（簡易版）
func (s *TestStorage) CreateTestMiniCommit(t *testing.T, message, patch string) string {
	// 簡易的なID生成
	id := fmt.Sprintf("test-%d", time.Now().UnixNano())
	
	// patchファイルを作成
	patchPath := filepath.Join(s.BasePath, id+".patch")
	if err := os.WriteFile(patchPath, []byte(patch), 0644); err != nil {
		t.Fatalf("Failed to create patch file: %v", err)
	}
	
	return id
}

// AssertMiniCommitExists mini-commitが存在することを確認
func (s *TestStorage) AssertMiniCommitExists(t *testing.T, id string) {
	patchPath := filepath.Join(s.BasePath, id+".patch")
	if _, err := os.Stat(patchPath); os.IsNotExist(err) {
		t.Errorf("Expected mini-commit %s to exist, but patch file not found", id)
	}
}

// AssertMiniCommitNotExists mini-commitが存在しないことを確認
func (s *TestStorage) AssertMiniCommitNotExists(t *testing.T, id string) {
	patchPath := filepath.Join(s.BasePath, id+".patch")
	if _, err := os.Stat(patchPath); !os.IsNotExist(err) {
		t.Errorf("Expected mini-commit %s to not exist, but patch file found", id)
	}
}

// AssertMiniCommitCount mini-commitの数を確認
func (s *TestStorage) AssertMiniCommitCount(t *testing.T, expected int) {
	files, err := os.ReadDir(s.BasePath)
	if err != nil {
		t.Fatalf("Failed to read mini-commits directory: %v", err)
	}
	
	patchCount := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".patch") {
			patchCount++
		}
	}
	
	if patchCount != expected {
		t.Errorf("Expected %d mini-commits, but got %d", expected, patchCount)
	}
}

// TestCLI テスト用CLI実行
type TestCLI struct{}

// NewTestCLI テスト用CLIを作成
func NewTestCLI(t *testing.T) *TestCLI {
	return &TestCLI{}
}

// RunCommand CLIコマンドを実行
func (c *TestCLI) RunCommand(args ...string) (string, string, error) {
	// 元のプロジェクトディレクトリのバイナリを使用
	// 環境変数から元のディレクトリを取得するか、固定パスを使用
	projectDir := os.Getenv("GIT_MINI_COMMIT_PROJECT_DIR")
	if projectDir == "" {
		// 環境変数が設定されていない場合は、現在のディレクトリから遡って探す
		wd, err := os.Getwd()
		if err != nil {
			return "", "", err
		}
		
		// 現在のディレクトリから遡ってgit-mini-commitバイナリを探す
		for {
			binaryPath := filepath.Join(wd, "git-mini-commit")
			if _, err := os.Stat(binaryPath); err == nil {
				projectDir = wd
				break
			}
			
			parent := filepath.Dir(wd)
			if parent == wd {
				break
			}
			wd = parent
		}
	}
	
	if projectDir == "" {
		return "", "", fmt.Errorf("git-mini-commit binary not found")
	}
	
	binaryPath := filepath.Join(projectDir, "git-mini-commit")
	
	cmd := exec.Command(binaryPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	
	return stdout.String(), stderr.String(), err
}

// AssertCommandSuccess コマンドが成功することを確認
func (c *TestCLI) AssertCommandSuccess(t *testing.T, args ...string) string {
	stdout, stderr, err := c.RunCommand(args...)
	if err != nil {
		t.Errorf("Command failed: %v, stderr: %s", err, stderr)
	}
	return stdout
}

// AssertCommandFailure コマンドが失敗することを確認
func (c *TestCLI) AssertCommandFailure(t *testing.T, args ...string) string {
	stdout, stderr, err := c.RunCommand(args...)
	if err == nil {
		t.Errorf("Expected command to fail, but it succeeded. stdout: %s", stdout)
	}
	// stdoutとstderrの両方を結合して返す
	combined := stdout + stderr
	return combined
}

// AssertOutputContains 出力に特定の文字列が含まれることを確認
func (c *TestCLI) AssertOutputContains(t *testing.T, output, expected string) {
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain '%s', but got: %s", expected, output)
	}
}

// AssertOutputNotContains 出力に特定の文字列が含まれないことを確認
func (c *TestCLI) AssertOutputNotContains(t *testing.T, output, unexpected string) {
	if strings.Contains(output, unexpected) {
		t.Errorf("Expected output to not contain '%s', but got: %s", unexpected, output)
	}
}

// TestFile テスト用ファイル管理
type TestFile struct {
	Path string
	Content string
}

// NewTestFile テスト用ファイルを作成
func NewTestFile(t *testing.T, path, content string) *TestFile {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file %s: %v", path, err)
	}
	
	return &TestFile{
		Path: path,
		Content: content,
	}
}

// Modify ファイルを修正
func (f *TestFile) Modify(t *testing.T, newContent string) {
	f.Content = newContent
	if err := os.WriteFile(f.Path, []byte(newContent), 0644); err != nil {
		t.Fatalf("Failed to modify test file %s: %v", f.Path, err)
	}
}

// AssertContent ファイルの内容を確認
func (f *TestFile) AssertContent(t *testing.T, expected string) {
	content, err := os.ReadFile(f.Path)
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", f.Path, err)
	}
	
	if string(content) != expected {
		t.Errorf("Expected file content to be '%s', but got '%s'", expected, string(content))
	}
}

// TestPatch テスト用patch管理
type TestPatch struct {
	Content string
}

// NewTestPatch テスト用patchを作成
func NewTestPatch(t *testing.T, content string) *TestPatch {
	return &TestPatch{
		Content: content,
	}
}

// ApplyToStaging patchをステージングエリアに適用
func (p *TestPatch) ApplyToStaging() error {
	cmd := exec.Command("git", "apply", "--cached")
	cmd.Stdin = strings.NewReader(p.Content)
	return cmd.Run()
}

// AssertPatchContent patchの内容を確認
func (p *TestPatch) AssertPatchContent(t *testing.T, expected string) {
	if p.Content != expected {
		t.Errorf("Expected patch content to be '%s', but got '%s'", expected, p.Content)
	}
}

// TestTimer テスト用タイマー
type TestTimer struct {
	FixedTime time.Time
}

// NewTestTimer テスト用タイマーを作成
func NewTestTimer(t *testing.T, fixedTime time.Time) *TestTimer {
	return &TestTimer{
		FixedTime: fixedTime,
	}
}

// Now 固定時刻を返す
func (tm *TestTimer) Now() time.Time {
	return tm.FixedTime
}

// TestLogger テスト用ロガー
type TestLogger struct {
	Logs []string
}

// NewTestLogger テスト用ロガーを作成
func NewTestLogger(t *testing.T) *TestLogger {
	return &TestLogger{
		Logs: make([]string, 0),
	}
}

// Log ログを記録
func (l *TestLogger) Log(message string) {
	l.Logs = append(l.Logs, message)
}

// AssertLogContains ログに特定の文字列が含まれることを確認
func (l *TestLogger) AssertLogContains(t *testing.T, expected string) {
	found := false
	for _, log := range l.Logs {
		if strings.Contains(log, expected) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected log to contain '%s', but got: %v", expected, l.Logs)
	}
}

// AssertLogCount ログの数を確認
func (l *TestLogger) AssertLogCount(t *testing.T, expected int) {
	if len(l.Logs) != expected {
		t.Errorf("Expected %d logs, but got %d: %v", expected, len(l.Logs), l.Logs)
	}
}
