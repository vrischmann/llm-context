package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

var version = "dev"

func main() {
	// Parse flags
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("llm-context version %s\n", version)
		os.Exit(0)
	}
	// 1. Check for fzf dependency
	if _, err := exec.LookPath("fzf"); err != nil {
		fmt.Println("Error: fzf is not installed or not in your PATH")
		os.Exit(1)
	}

	// 2. Generate File List
	// We try to use 'fd' or 'git' first for better ignoring, fallback to standard directory walk
	files, err := getFileList()
	if err != nil {
		fmt.Printf("Error listing files: %v\n", err)
		os.Exit(1)
	}

	// 3. Run FZF
	selectedFiles, err := runFzf(files)
	if err != nil {
		// If prompt was cancelled (exit code 130), just exit quietly
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 130 {
			os.Exit(0)
		}
		fmt.Printf("Error running fzf: %v\n", err)
		os.Exit(1)
	}

	if len(selectedFiles) == 0 {
		fmt.Println("No files selected.")
		os.Exit(0)
	}

	// 4. Build Context
	var sb strings.Builder
	sb.WriteString("I am providing the following files as context:\n\n")

	for _, file := range selectedFiles {
		if file == "" {
			continue
		}
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Warning: Could not read file %s: %v\n", file, err)
			continue
		}

		ext := strings.TrimPrefix(filepath.Ext(file), ".")
		if ext == "" {
			ext = "text"
		}

		sb.WriteString(fmt.Sprintf("## File: %s\n", file))
		sb.WriteString(fmt.Sprintf("```%s\n", ext))
		sb.WriteString(string(content))
		sb.WriteString("\n```\n\n")
	}

	// 5. Copy to Clipboard
	finalOutput := sb.String()
	err = clipboard.WriteAll(finalOutput)
	if err != nil {
		fmt.Printf("⚠️  Failed to copy to clipboard: %v\n", err)
		fmt.Println("Here is the output instead:")
		fmt.Println("---------------------------")
		fmt.Println(finalOutput)
	} else {
		fmt.Printf("✅ Copied %d files to clipboard (%d chars)!\n", len(selectedFiles), len(finalOutput))
	}
}

// getFileList tries to be smart about how it finds files
func getFileList() (string, error) {
	// Strategy 1: Try 'fd' (fast, respects gitignore)
	if _, err := exec.LookPath("fd"); err == nil {
		cmd := exec.Command("fd", "--type", "f")
		out, err := cmd.Output()
		if err == nil {
			return string(out), nil
		}
	}

	// Strategy 2: Try 'git ls-files' (respects gitignore)
	if _, err := exec.LookPath("git"); err == nil {
		// Check if we are actually in a git repo
		if err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Run(); err == nil {
			cmd := exec.Command("git", "ls-files")
			out, err := cmd.Output()
			if err == nil {
				return string(out), nil
			}
		}
	}

	// Strategy 3: Native Go Walk (Fallback)
	var files []string
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip hidden directories like .git
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return strings.Join(files, "\n"), err
}

func runFzf(input string) ([]string, error) {
	cmd := exec.Command("fzf", "--multi", "--height=80%", "--layout=reverse", "--border", "--preview=head -n 20 {}")

	// Create a pipe to stdin of fzf
	cmd.Stdin = strings.NewReader(input)

	// FZF UI needs to print to Stderr to be visible to the user
	cmd.Stderr = os.Stderr

	// Capture the output (selected files)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// Split output by newlines
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	return lines, nil
}
