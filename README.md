# Kommito

Kommito is a lightweight version control system inspired by Git, built as a fun learning project to understand the internals of version control systems. It implements core Git concepts in a simplified way while maintaining a clean, modular architecture.

> 🎮 **Personal Project**: This is a fun side project I built to learn about version control systems. Just sharing my learning journey!

> ⚠️ **Development Status**: This project is in active development as I learn and experiment. Features are added as I explore different aspects of version control systems.

## Why Kommito?

I built Kommito to:
- Learn and understand the fundamental concepts of version control systems
- Implement Git-like functionality in a simplified manner
- Have fun experimenting with file system operations and data structures
- Challenge myself to build something complex from scratch

## Core Features

Kommito provides essential version control features:

1. **Repository Management**
   - Initialize new repositories
   - Track file changes
   - Maintain commit history
   - Manage file versions
   - Clone existing repositories (both Kommito and Git)

2. **File Operations**
   - Stage individual or multiple files
   - Track file modifications
   - Store file contents efficiently
   - Handle file conflicts
   - Skip system files automatically

3. **Commit System**
   - Create commits with messages
   - Track commit history
   - Maintain commit metadata
   - Link commits in a chain
   - Automatic author detection

4. **Status and History**
   - View repository status
   - Check staged changes
   - Display commit history
   - Show file modifications
   - Track untracked files

5. **Repository Cloning**
   - Clone local Kommito repositories
   - Clone Git repositories
   - Automatic conversion from Git to Kommito format
   - Preserve file history and structure

## Project Structure

```
kommito/
├── internal/           # Internal packages
│   └── repo/          # Core repository operations
│       ├── objects.go  # Object storage handling
│       ├── commit.go   # Commit operations
│       ├── index.go    # Staging area management
│       ├── clone.go    # Repository cloning
│       ├── log.go      # Commit history
│       └── status.go   # Repository status
├── main.go            # CLI entry point
└── go.mod             # Go module definition
```

## Repository Structure

When you initialize a Kommito repository, it creates the following structure:

```
.kommito/
├── objects/           # Object storage
│   ├── blobs/        # File contents (SHA-1 hashed)
│   └── commits/      # Commit objects (JSON format)
├── refs/             # References
│   └── heads/        # Branch references
├── HEAD              # Points to current commit
├── index            # Staging area (JSON format)
└── config.json      # Repository configuration
```

## Implementation Details

### Object Storage System

1. **Blob Objects**
   - Store actual file contents
   - Named using SHA-1 hash of content
   - Compressed for efficiency
   - Located in `.kommito/objects/blobs/`

2. **Commit Objects**
   - Store commit metadata
   - Reference parent commits
   - List of blob references
   - Located in `.kommito/objects/commits/`

### Data Structures

```go
// Commit structure
type Commit struct {
    Author    string   `json:"author"`     // Commit author
    Timestamp string   `json:"timestamp"`  // Commit timestamp
    Message   string   `json:"message"`    // Commit message
    Blobs     []string `json:"blobs"`      // Referenced blob hashes
}

// Index structure
type Index struct {
    Files map[string]string `json:"files"` // File path to blob hash mapping
}
```

## Usage Guide

### Basic Commands

```bash
# Initialize a new repository
kommito init

# Stage files
kommito add <file>    # Stage a single file
kommito add .         # Stage all files in current directory

# Create a commit
kommito commit -m "Your commit message"

# View commit history
kommito log

# Check repository status
kommito status

# Clone a repository
kommito clone <source> <destination>  # Clone local Kommito repo
kommito clone <git-url> <destination> # Clone Git repo
```

### Workflow Examples

1. **Start a New Project**
   ```bash
   mkdir my-project
   cd my-project
   kommito init
   ```

2. **Add and Track Files**
   ```bash
   echo "# My Project" > README.md
   kommito add README.md
   kommito commit -m "Initial commit"
   ```

3. **Clone an Existing Repository**
   ```bash
   # Clone a local Kommito repository
   kommito clone /path/to/repo my-clone

   # Clone a Git repository
   kommito clone https://github.com/user/repo.git my-clone
   ```

4. **Make Changes and Commit**
   ```bash
   # Edit files
   kommito add .
   kommito commit -m "Update project files"
   ```

## Technical Implementation

### Core Technologies

- **Language**: Go (Golang)
- **File Operations**: Go standard library (`os`, `io`)
- **Hashing**: SHA-1 for content addressing
- **Data Storage**: JSON for metadata
- **Error Handling**: Custom error types and messages
- **CLI Framework**: Cobra for command-line interface

### Key Features

1. **Atomic Operations**
   - Safe file writes
   - Transaction-like commits
   - Error recovery

2. **Efficient Storage**
   - Content-based addressing
   - Deduplication of content
   - Compressed storage

3. **User Experience**
   - Simple command interface
   - Clear error messages
   - Intuitive workflow
   - Emoji-based status indicators

4. **Git Compatibility**
   - Clone Git repositories
   - Convert Git history
   - Preserve file structure

## Building and Installation

### Prerequisites
- Go 1.16 or higher
- Git (for version control and cloning)

```bash
# Clone the repository
git clone https://github.com/Kshitijknk07/Kommito.git
cd Kommito

# Build the project
go build

# Run directly
go run main.go
```

### Important Notes

1. **Current Limitations**
   - Basic version control features only
   - No remote repository support
   - Limited branch management
   - No merge conflict resolution
   - Built for learning purposes only

## License

This project is open source and available under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by Git's design and implementation
- Built as a personal learning project
- Thanks to the Go community for excellent documentation and tools
