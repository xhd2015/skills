package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func handleYAML(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("yaml requires a subcommand: validate")
	}

	switch args[0] {
	case "validate":
		return handleYAMLValidate(args[1:])
	case "-h", "--help":
		fmt.Print(yamlHelp)
		return nil
	default:
		return fmt.Errorf("unknown yaml subcommand: %s (expected `validate`)", args[0])
	}
}

func handleYAMLValidate(args []string) error {
	var filePath string
	for _, a := range args {
		if a == "-h" || a == "--help" {
			fmt.Print(yamlValidateHelp)
			return nil
		}
		if !strings.HasPrefix(a, "-") {
			filePath = a
		}
	}

	if filePath == "" {
		return fmt.Errorf("yaml validate requires a file path")
	}

	return validateWorkflowFile(filePath)
}

func isActionlintAvailable() bool {
	_, err := exec.LookPath("actionlint")
	return err == nil
}

func validateWorkflowFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		}
		return fmt.Errorf("read file %s: %w", filePath, err)
	}

	if isActionlintAvailable() {
		return runActionlint(filePath)
	}

	return localValidateYAML(filePath, string(data))
}

func runActionlint(filePath string) error {
	cmd := exec.Command("actionlint", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var (
	reOn   = regexp.MustCompile(`(?m)^\s*on\s*:`)
	reJobs = regexp.MustCompile(`(?m)^\s*jobs\s*:`)
)

func localValidateYAML(filePath, content string) error {
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(content), &node); err != nil {
		fmt.Printf("%s: YAML syntax error\n  %v\n", filePath, err)
		return nil
	}

	var issues []string
	if !reOn.MatchString(content) {
		issues = append(issues, "missing `on:` trigger definition")
	}
	if !reJobs.MatchString(content) {
		issues = append(issues, "missing `jobs:` definition")
	}

	if len(issues) == 0 {
		fmt.Printf("%s: OK (valid YAML, basic checks passed)\n", filePath)
		fmt.Println("(install actionlint for full GitHub Actions schema validation: brew install actionlint)")
		return nil
	}

	fmt.Printf("%s: validation issues\n", filePath)
	for _, issue := range issues {
		fmt.Printf("  - %s\n", issue)
	}
	fmt.Println("(install actionlint for full GitHub Actions schema validation: brew install actionlint)")
	return nil
}

const yamlHelp = `
Usage: github-fetch yaml <subcommand> [ARGS]

Subcommands:
  validate <path>   Validate a GitHub Actions workflow file

Options:
  -h, --help        Show this help message
`

const yamlValidateHelp = `
Usage: github-fetch yaml validate <path>

Validate a GitHub Actions workflow YAML file.

Two-level validation:
  1. YAML syntax check (parse errors, indentation, etc.)
  2. GitHub Actions semantic checks (on:, jobs: required fields)

If actionlint (https://github.com/rhysd/actionlint) is installed,
it is used for comprehensive schema validation instead.
`
