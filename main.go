package main

import (
	"fmt"
	"os"

	repo "github.com/Kshitijknk07/Kommito/internal/repo"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println(`(｡•́︿•̀｡) Nani?! You forgot the command!

🔮 Usage:
   kommito <command>

🧭 Available Commands:
   init    ⚙️  Initialize a brand new Kommito repo
   add     ➕  Stage files for commit
   commit  📝  Commit staged files
   log     📜  Show commit history
   status  ��  Show repo status
   clone   📋  Clone a repository

✨ Example:
   kommito init
   kommito add <file>    # Stage a single file
   kommito add .         # Stage all files
   kommito commit -m "message"
   kommito log
   kommito status
   kommito clone <source> <destination>`)
		return
	}

	switch args[1] {
	case "init":
		fmt.Println(`(｀• ω •´)ゞ Roger that!

⚒️  Spinning up your Kommito engine...
🗂️  Setting up the repository chamber...`)

		if err := repo.InitRepo(); err != nil {
			fmt.Printf("(╥﹏╥) Oops! Something went wrong: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✨ Repository initialized successfully!")
	case "add":
		if len(args) < 3 {
			fmt.Println(`(⊙_☉) You need to specify a file to add!

✨ Example:
   kommito add myfile.txt    # Stage a single file
   kommito add .             # Stage all files`)
			return
		}
		filePath := args[2]
		fmt.Printf("(ง •_•)ง Staging files...\n")
		if err := repo.AddFile(filePath); err != nil {
			fmt.Printf("(╥﹏╥) Could not add files: %v\n", err)
			os.Exit(1)
		}
	case "commit":
		if len(args) < 4 || args[2] != "-m" {
			fmt.Println(`(⊙_☉) You need to provide a commit message!

✨ Example:
   kommito commit -m "Initial commit"`)
			return
		}
		message := args[3]
		fmt.Println("(ﾉ◕ヮ◕)ﾉ*:･ﾟ✧ Creating your commit...")
		if err := repo.CommitStaged(message); err != nil {
			fmt.Printf("(╥﹏╥) Commit failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("(づ｡◕‿‿◕｡)づ Commit created successfully!")
	case "log":
		if err := repo.LogCommits(); err != nil {
			fmt.Printf("(╥﹏╥) Could not show log: %v\n", err)
		}
	case "status":
		if err := repo.Status(); err != nil {
			fmt.Printf("(╥﹏╥) Could not show status: %v\n", err)
		}
	case "clone":
		if len(args) < 4 {
			fmt.Println(`(⊙_☉) You need to specify source and destination!

✨ Example:
   kommito clone /path/to/source /path/to/destination`)
			return
		}
		source := args[2]
		destination := args[3]
		fmt.Println("(ﾉ◕ヮ◕)ﾉ*:･ﾟ✧ Cloning repository...")
		if err := repo.CloneRepo(source, destination); err != nil {
			fmt.Printf("(╥﹏╥) Clone failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("(づ｡◕‿‿◕｡)づ Repository cloned successfully!")
	default:
		fmt.Printf(`(¬_¬) I don't know that command: "%s"

Maybe try:
   kommito init
   kommito add <file>    # Stage a single file
   kommito add .         # Stage all files
   kommito commit -m "message"
   kommito log
   kommito status
   kommito clone <source> <destination>

Kommito is still just a chibi tool... be nice to it! 🐣`, args[1])
	}
}
