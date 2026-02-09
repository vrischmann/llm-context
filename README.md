# llm-context

A command-line tool that helps you quickly gather file contents into a formatted message that can be pasted into chatbot conversations. Perfect for providing code context to LLMs like GitHub Copilot, Claude, or other AI assistants.

## Features

- **Interactive file selection**: Uses fzf for easy multi-file selection with preview
- **Smart file discovery**: Automatically finds files using fd, git ls-files, or native Go walk
- **Formatted output**: Creates Markdown-formatted context with file names and syntax highlighting
- **Clipboard integration**: Copies the formatted context directly to your clipboard
- **Git-aware**: Respects .gitignore files when available

## How It Works

1. Scans your current directory for files (respecting .gitignore if in a git repo)
2. Presents an interactive fzf selector with file previews
3. Lets you select multiple files
4. Formats the selected files into a Markdown message with proper code blocks
5. Copies the result to your clipboard for easy pasting into chat interfaces

## Requirements

- Go 1.20+
- [fzf](https://github.com/junegunn/fzf) (for interactive file selection)
- Optional but recommended: [fd](https://github.com/sharkdp/fd) (for faster file discovery)

## Installation

```bash
git clone https://github.com/yourusername/llm-context.git
cd llm-context
go build -o llm-context
```

## Usage

```bash
# Run from any directory
./llm-context

# Or run directly with Go
go run main.go
```

## Example Output Format

The tool generates output like this:

```markdown
I am providing the following files as context:

## File: main.go
```go
package main

func main() {
    // your code here
}
```

## File: config.yaml
```yaml
key: value
```
```

## License

MIT
