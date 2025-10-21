package cmd

import (
	"fmt"
	"os"
	"time"

	"git-mini-commit/internal/git"
	"git-mini-commit/internal/storage"
	"git-mini-commit/internal/types"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-mini-commit",
	Short: "Manage mini-commits between staging area and regular commits",
	Long: `git-mini-commit is a tool that introduces "mini-commit" as an intermediate unit between Git's staging area and regular commits, making it easier to manage diffs during large-scale refactoring or change work.

The staging area content is saved in small units and managed locally only.

Usage:
  git mini-commit -m "message"      # Create mini-commit
  git mini-commit list              # List mini-commits
  git mini-commit show <hash>       # Show mini-commit diff
  git mini-commit pop <hash>        # Apply mini-commit to staging
  git mini-commit drop <hash>       # Delete mini-commit
  git commit -m "message"          # Integration commit (standard Git command)`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		message, _ := cmd.Flags().GetString("message")
		if message == "" {
			return fmt.Errorf("message is required (-m option)")
		}

		// Gitリポジトリかチェック
		if !git.IsGitRepository() {
			return fmt.Errorf("not a git repository")
		}

		// ステージングされた変更があるかチェック
		hasChanges, err := git.HasStagedChanges()
		if err != nil {
			return fmt.Errorf("failed to check staging status: %v", err)
		}
		if !hasChanges {
			return fmt.Errorf("no staged changes")
		}

		// ステージングされた変更を取得
		patch, err := git.GetStagedChanges()
		if err != nil {
			return fmt.Errorf("failed to get staged changes: %v", err)
		}

		// ストレージを初期化
		storage, err := storage.NewStorage()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %v", err)
		}

		// mini-commitを作成
		now := time.Now()
		mc := &types.MiniCommit{
			ID:        storage.GenerateID(patch, now),
			Message:   message,
			CreatedAt: now,
			Patch:     patch,
		}

		// 保存
		if err := storage.SaveMiniCommit(mc); err != nil {
			return fmt.Errorf("failed to save mini-commit: %v", err)
		}

		fmt.Printf("Created mini-commit: %s\n", mc.ID[:8])
		fmt.Printf("Message: %s\n", mc.Message)
		fmt.Printf("Created at: %s\n", mc.CreatedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

func init() {
	rootCmd.Flags().StringP("message", "m", "", "mini-commit message")
}

// Execute はコマンドを実行します
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
