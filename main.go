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
   add     ➕  Stage a file for commit

✨ Example:
   kommito init
   kommito add <file>`)
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
   kommito add myfile.txt`)
			return
		}
		filePath := args[2]
		fmt.Printf("(ง •_•)ง Staging file: %s ...\n", filePath)
		if err := repo.AddFile(filePath); err != nil {
			fmt.Printf("(╥﹏╥) Could not add file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("(＾▽＾) File staged successfully!")
	default:
		fmt.Printf(`(¬_¬) I don't know that command: "%s"

Maybe try:
   kommito init
   kommito add <file>

Kommito is still just a chibi tool... be nice to it! 🐣`, args[1])
	}
}
