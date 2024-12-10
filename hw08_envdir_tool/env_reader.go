package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.Contains(name, "=") {
			return nil, fmt.Errorf("invalid environment variable name: %s", name)
		}

		filePath := filepath.Join(dir, name)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("could not read file %s: %w", filePath, err)
		}

		// Удаляем нулевые байты и пробелы в конце строки
		value := strings.TrimRight(string(content), " \t\n\x00")
		env[name] = EnvValue{
			Value:      value,
			NeedRemove: len(value) == 0,
		}
	}

	return env, nil
}
