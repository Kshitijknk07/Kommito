# Kommito

Kommito is a lightweight version control system inspired by Git, built as a fun learning project to understand the internals of version control systems. It implements core Git concepts in a simplified way while maintaining a clean, modular architecture.

## Overview

Kommito provides essential version control features:
- Repository initialization
- File staging
- Commit creation
- Commit history viewing
- Repository status checking

## Project Structure

```
kommito/
├── internal/      # Internal packages
│   └── repo/     # Core repository operations
├── main.go       # Entry point
└── go.mod        # Go module definition
```

## Core Components

### Repository Structure
```
.kommito/
├── objects/
│   ├── blobs/    # File contents
│   └── commits/  # Commit objects
├── refs/
│   └── heads/    # Branch references
├── HEAD          # Current commit reference
├── index         # Staging area
└── config.json   # Repository configuration
```

### Key Features

1. **Repository Initialization**
   - Creates necessary directory structure
   - Initializes configuration
   - Sets up HEAD and index

2. **File Staging**
   - Single file staging: `kommito add <file>`
   - Batch staging: `kommito add .`
   - Creates blob objects
   - Updates index

3. **Commit Creation**
   - Captures staged files
   - Records author and timestamp
   - Creates commit object
   - Updates HEAD reference

4. **Status Checking**
   - Shows staged files
   - Lists modified files
   - Identifies untracked files

5. **Commit History**
   - Displays commit information
   - Shows author and timestamp
   - Lists commit message

## Implementation Details

### Object Storage
- Files are stored as blobs with SHA-1 hashes
- Commits are stored as JSON objects
- Index maintains file-to-blob mapping

### Data Structures
```go
type Commit struct {
    Author    string   `json:"author"`
    Timestamp string   `json:"timestamp"`
    Message   string   `json:"message"`
    Blobs     []string `json:"blobs"`
}
```

## Usage

```bash
# Initialize repository
kommito init

# Stage files
kommito add <file>    # Stage single file
kommito add .         # Stage all files

# Create commit
kommito commit -m "message"

# View history
kommito log

# Check status
kommito status
```

## Technical Implementation

### File Operations
- Uses Go's standard library for file operations
- Implements SHA-1 hashing for content addressing
- Maintains atomic file operations

### Error Handling
- Comprehensive error checking
- User-friendly error messages
- Graceful failure handling

### Code Organization
- Clean separation of concerns
- Modular package structure
- Clear function documentation

## Building and Installation

```bash
# Build and install
go install

# Run directly
go run main.go
```

## Acknowledgments

This project was created as a learning exercise to understand version control systems. It draws inspiration from Git's design while maintaining simplicity and clarity.

## License

This project is open source and available under the MIT License.
