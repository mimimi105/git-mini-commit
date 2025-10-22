package cmd

import (
	"fmt"

	"git-mini-commit/internal/git"
	"git-mini-commit/internal/storage"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <hash>",
	Short: "Show diff of specified mini-commit",
	Long:  `Display the diff (patch) of the mini-commit with the specified ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash := args[0]

		// Check if it's a Git repository
		if !git.IsGitRepository() {
			return fmt.Errorf("not a git repository")
		}

		// Initialize storage
		storage, err := storage.NewStorage()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %v", err)
		}

		// Get mini-commit
		mc, err := storage.GetMiniCommit(hash)
		if err != nil {
			return fmt.Errorf("failed to get mini-commit: %v", err)
		}

		// Display information
		fmt.Printf("Mini-commit: %s\n", mc.ID[:8])
		fmt.Printf("Message: %s\n", mc.Message)
		fmt.Printf("Created: %s\n", mc.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("\nDiff:")
		fmt.Println("---")
		fmt.Print(mc.Patch)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
