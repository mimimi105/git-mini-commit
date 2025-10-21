package cmd

import (
	"fmt"

	"git-mini-commit/internal/git"
	"git-mini-commit/internal/storage"

	"github.com/spf13/cobra"
)

var dropCmd = &cobra.Command{
	Use:   "drop <hash>",
	Short: "Delete specified mini-commit",
	Long:  `Delete the mini-commit with the specified ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash := args[0]

		// Gitリポジトリかチェック
		if !git.IsGitRepository() {
			return fmt.Errorf("not a git repository")
		}

		// ストレージを初期化
		storage, err := storage.NewStorage()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %v", err)
		}

		// mini-commitを削除
		if err := storage.DeleteMiniCommit(hash); err != nil {
			return fmt.Errorf("failed to delete mini-commit: %v", err)
		}

		fmt.Printf("Deleted mini-commit '%s'\n", hash[:8])

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dropCmd)
}
