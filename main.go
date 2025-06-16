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
   clone   📋  Clone a repository`,
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
	Use:   "commit -m [message]",
	Short: "Commit staged files",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "-m" {
			fmt.Println(`(⊙_☉) You need to provide a commit message!

✨ Example:
   kommito commit -m "Initial commit"`)
			return
		}
		message := args[1]
		fmt.Println("(ﾉ◕ヮ◕)ﾉ*:･ﾟ✧ Creating your commit...")
		if err := repo.CommitStaged(message); err != nil {
			fmt.Printf("(╥﹏╥) Commit failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("(づ｡◕‿‿◕｡)づ Commit created successfully!")
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

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(cloneCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
