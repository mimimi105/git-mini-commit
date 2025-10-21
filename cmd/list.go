package cmd

import (
	"fmt"

	"git-mini-commit/internal/git"
	"git-mini-commit/internal/storage"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved mini-commits",
	Long:  `Display a list of all saved mini-commits.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Gitリポジトリかチェック
		if !git.IsGitRepository() {
			return fmt.Errorf("not a git repository")
		}

		// ストレージを初期化
		storage, err := storage.NewStorage()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %v", err)
		}

		// mini-commit一覧を取得
		miniCommits, err := storage.LoadMiniCommits()
		if err != nil {
			return fmt.Errorf("failed to load mini-commits: %v", err)
		}

		if len(miniCommits) == 0 {
			fmt.Println("No mini-commits found")
			return nil
		}

		// 一覧を表示
		fmt.Printf("Mini-commits (%d):\n\n", len(miniCommits))
		for i, mc := range miniCommits {
			fmt.Printf("%d. ID: %s\n", i+1, mc.ID[:8])
			fmt.Printf("   Message: %s\n", mc.Message)
			fmt.Printf("   Created: %s\n", mc.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
