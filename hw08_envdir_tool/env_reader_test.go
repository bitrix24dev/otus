package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("read environment variables from directory", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "envdir")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		err = os.WriteFile(filepath.Join(dir, "FOO"), []byte("123\n"), 0o644)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(dir, "BAR"), []byte("value\n"), 0o644)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(dir, "EMPTY"), []byte(""), 0o644)
		require.NoError(t, err)

		env, err := ReadDir(dir)
		require.NoError(t, err)

		require.Equal(t, Environment{
			"FOO":   {Value: "123", NeedRemove: false},
			"BAR":   {Value: "value", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: true},
		}, env)
	})
}
