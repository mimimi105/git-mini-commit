package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Git操作のユーティリティ関数群

// GetStagedChanges ステージングされた変更をpatch形式で取得
func GetStagedChanges() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get staged changes: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// HasStagedChanges ステージングされた変更があるかチェック
func HasStagedChanges() (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()

	// exit code 0: 変更なし, 1: 変更あり
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return true, nil
		}
		return false, fmt.Errorf("failed to check staging status: %v", err)
	}

	return false, nil
}

// ApplyPatch patchをステージングエリアに適用
func ApplyPatch(patch string) error {
	cmd := exec.Command("git", "apply", "--cached")
	cmd.Stdin = strings.NewReader(patch)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to apply patch: %v, stderr: %s", err, stderr.String())
	}

	return nil
}

// GetWorkingDirectoryChanges ワーキングディレクトリの変更をpatch形式で取得
func GetWorkingDirectoryChanges() (string, error) {
	cmd := exec.Command("git", "diff")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get working directory changes: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// CommitWithMessage 指定されたメッセージでコミット
func CommitWithMessage(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %v, stderr: %s", err, stderr.String())
	}

	return nil
}

// IsGitRepository 現在のディレクトリがGitリポジトリかチェック
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// GetRepositoryRoot Gitリポジトリのルートディレクトリを取得
func GetRepositoryRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get repository root: %v, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}
