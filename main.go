package main

import (
	"fmt"
	"os"

	repo "github.com/Kshitijknk07/Kommito/internal/repo"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kommito",
	Short: "Kommito - A lightweight version control system",
	Long: `(｡•́︿•̀｡) Kommito is a lightweight version control system inspired by Git.

🔮 Usage:
   kommito <command>

🧭 Available Commands:
   init    ⚙️  Initialize a brand new Kommito repo
   add     ➕  Stage files for commit
   commit  📝  Commit staged files
   log     📜  Show commit history
   status  🧭  Show repo status
   clone   📋  Clone a repository
   branch  🌿  Manage branches`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Kommito repository",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`(｀• ω •´)ゞ Roger that!

⚒️  Spinning up your Kommito engine...
🗂️  Setting up the repository chamber...`)

		if err := repo.InitRepo(); err != nil {
			fmt.Printf("(╥﹏╥) Oops! Something went wrong: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✨ Repository initialized successfully!")
	},
}

var addCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Stage files for commit",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		fmt.Printf("(ง •_•)ง Staging files...\n")
		if err := repo.AddFile(filePath); err != nil {
			fmt.Printf("(╥﹏╥) Could not add files: %v\n", err)
			os.Exit(1)
		}
	},
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit staged files",
	RunE: func(cmd *cobra.Command, args []string) error {
		message, _ := cmd.Flags().GetString("message")
		if message == "" {
			return fmt.Errorf(`(⊙_☉) You need to provide a commit message!

✨ Example:
   kommito commit --message "Initial commit"`)
		}
		fmt.Println("(ﾉ◕ヮ◕)ﾉ*:･ﾟ✧ Creating your commit...")
		if err := repo.CommitStaged(message); err != nil {
			return fmt.Errorf("(╥﹏╥) Commit failed: %v", err)
		}
		fmt.Println("(づ｡◕‿‿◕｡)づ Commit created successfully!")
		return nil
	},
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit history",
	Run: func(cmd *cobra.Command, args []string) {
		if err := repo.LogCommits(); err != nil {
			fmt.Printf("(╥﹏╥) Could not show log: %v\n", err)
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show repository status",
	Run: func(cmd *cobra.Command, args []string) {
		if err := repo.Status(); err != nil {
			fmt.Printf("(╥﹏╥) Could not show status: %v\n", err)
		}
	},
}

var cloneCmd = &cobra.Command{
	Use:   "clone [source] [destination]",
	Short: "Clone a repository",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]
		destination := args[1]
		fmt.Println("(ﾉ◕ヮ◕)ﾉ*:･ﾟ✧ Cloning repository...")
		if err := repo.CloneRepo(source, destination); err != nil {
			fmt.Printf("(╥﹏╥) Clone failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("(づ｡◕‿‿◕｡)づ Repository cloned successfully!")
	},
}

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Manage branches",
	Long: `Manage branches in your repository.

Available subcommands:
  list     List all branches
  create   Create a new branch
  switch   Switch to a branch
  delete   Delete a branch`,
}

var branchListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all branches",
	Run: func(cmd *cobra.Command, args []string) {
		bm := repo.NewBranchManager(".")
		branches, err := bm.ListBranches()
		if err != nil {
			fmt.Printf("(╥﹏╥) Could not list branches: %v\n", err)
			os.Exit(1)
		}

		currentBranch, _ := bm.GetCurrentBranch()
		fmt.Println("🌿 Branches:")
		for _, branch := range branches {
			marker := "  "
			if branch.Name == currentBranch {
				marker = "→ "
			}
			fmt.Printf("%s%s\n", marker, branch.Name)
		}
	},
}

var branchCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		bm := repo.NewBranchManager(".")
		if err := bm.CreateBranch(name); err != nil {
			fmt.Printf("(╥﹏╥) Could not create branch: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✨ Branch '%s' created successfully!\n", name)
	},
}

var branchSwitchCmd = &cobra.Command{
	Use:   "switch [name]",
	Short: "Switch to a branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		bm := repo.NewBranchManager(".")
		if err := bm.SwitchBranch(name); err != nil {
			fmt.Printf("(╥﹏╥) Could not switch branch: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✨ Switched to branch '%s'\n", name)
	},
}

var branchDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		bm := repo.NewBranchManager(".")
		if err := bm.DeleteBranch(name); err != nil {
			fmt.Printf("(╥﹏╥) Could not delete branch: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✨ Branch '%s' deleted successfully!\n", name)
	},
}

var mergeCmd = &cobra.Command{
	Use:   "merge [branch]",
	Short: "Merge a branch into the current branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branch := args[0]
		if err := repo.MergeBranches(branch); err != nil {
			fmt.Printf("Merge failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(cloneCmd)

	rootCmd.AddCommand(branchCmd)
	branchCmd.AddCommand(branchListCmd)
	branchCmd.AddCommand(branchCreateCmd)
	branchCmd.AddCommand(branchSwitchCmd)
	branchCmd.AddCommand(branchDeleteCmd)
	rootCmd.AddCommand(mergeCmd)

	commitCmd.Flags().StringP("message", "m", "", "Commit message")
	commitCmd.MarkFlagRequired("message")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
