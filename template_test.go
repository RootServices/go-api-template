package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestTemplateGenerationAndTest(t *testing.T) {
	// Find cookiecutter binary
	cookiecutterPath, err := exec.LookPath("cookiecutter")
	if err != nil {
		// Try common locations if not in PATH
		homeDir, _ := os.UserHomeDir()
		possiblePaths := []string{
			filepath.Join(homeDir, "Library/Python/3.9/bin/cookiecutter"),
			filepath.Join(homeDir, ".local/bin/cookiecutter"),
		}
		for _, p := range possiblePaths {
			if _, err := os.Stat(p); err == nil {
				cookiecutterPath = p
				break
			}
		}
	}

	if cookiecutterPath == "" {
		t.Skip("cookiecutter not found, skipping template test")
	}

	// Create a temporary directory for the generated project
	tempDir, err := os.MkdirTemp("", "go-api-template-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	// defer os.RemoveAll(tempDir)

	// Get the absolute path to the template directory (current directory)
	templateDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	t.Logf("Template directory: %s", templateDir)
	t.Logf("tempDir: %s", tempDir)
	// Generate the project
	cmd := exec.Command(cookiecutterPath, templateDir, "--no-input", "--output-dir", tempDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run cookiecutter: %v\nOutput: %s", err, output)
	}

	// The generated project name is "go-api-project" (from default in cookiecutter.json)
	projectDir := filepath.Join(tempDir, "go-api-project")

	// Verify the project directory exists
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatalf("generated project directory does not exist: %s", projectDir)
	}

	// Initialize git repository (required for Makefile version generation)
	cmd = exec.Command("git", "init")
	cmd.Dir = projectDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run git init: %v\nOutput: %s", err, output)
	}

	cmd = exec.Command("git", "config", "user.name", "tmack")
	cmd.Dir = projectDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git user.name: %v", err)
	}

	cmd = exec.Command("git", "config", "user.email", "tmack@fake-email.com")
	cmd.Dir = projectDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git user.email: %v", err)
	}

	// Create an initial commit so git rev-parse HEAD works
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = projectDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git add: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = projectDir

	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to add a git commit: %v", err)
	}

	// Run go mod tidy (since the post-gen hook was deleted)
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run go mod tidy: %v\nOutput: %s", err, output)
	}

	// Run make test
	cmd = exec.Command("make", "test")
	cmd.Dir = projectDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run make test: %v\nOutput: %s", err, output)
	}

	t.Logf("Successfully generated project and ran tests:\n%s", output)
}
