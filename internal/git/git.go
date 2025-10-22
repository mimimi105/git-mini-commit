package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Git operation utility functions

// GetStagedChanges gets staged changes in patch format
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

// HasStagedChanges checks if there are staged changes
func HasStagedChanges() (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()

	// exit code 0: no changes, 1: has changes
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return true, nil
		}
		return false, fmt.Errorf("failed to check staging status: %v", err)
	}

	return false, nil
}

// ApplyPatch applies patch to staging area
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


// IsGitRepository checks if current directory is a Git repository
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

