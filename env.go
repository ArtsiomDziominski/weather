// Package main provides a utility for loading environment variables from a
// .env file into the current process's environment.
package main

import (
	"bufio"
	"os"
	"strings"
)

// loadDotEnv reads a .env file and sets each KEY=VALUE pair as an environment
// variable for the current process. Lines starting with '#' and blank lines
// are ignored. Variables already set in the environment are NOT overwritten,
// so shell-level values always take precedence over the .env file.
// If the file does not exist, loadDotEnv returns nil (not an error).
func loadDotEnv(filename string) error {
	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil // .env is optional
	}
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on the first '=' only.
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue // malformed line — skip silently
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// Strip optional surrounding quotes from the value ("value" or 'value').
		if len(value) >= 2 &&
			((value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}

		// Don't overwrite variables already present in the environment.
		if os.Getenv(key) == "" {
			os.Setenv(key, value) //nolint:errcheck
		}
	}

	return scanner.Err()
}
