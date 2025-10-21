package cmd

import (
	"fmt"

	"git-mini-commit/internal/git"
	"git-mini-commit/internal/storage"

	"github.com/spf13/cobra"
)

var popCmd = &cobra.Command{
	Use:   "pop <hash>",
	Short: "Apply mini-commit content back to staging",
	Long:  `Apply the content of the mini-commit with the specified ID to the staging area.`,
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

		// mini-commitを取得
		mc, err := storage.GetMiniCommit(hash)
		if err != nil {
			return fmt.Errorf("failed to get mini-commit: %v", err)
		}

		// patchをステージングエリアに適用
		if err := git.ApplyPatch(mc.Patch); err != nil {
			return fmt.Errorf("failed to apply patch: %v", err)
		}

		fmt.Printf("Applied mini-commit '%s' to staging area\n", mc.ID[:8])
		fmt.Printf("Message: %s\n", mc.Message)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(popCmd)
}
