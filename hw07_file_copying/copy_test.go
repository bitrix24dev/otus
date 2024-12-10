package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("copy entire file", func(t *testing.T) {
		srcFile, err := os.CreateTemp("", "src")
		require.NoError(t, err)
		defer os.Remove(srcFile.Name())

		dstFile, err := os.CreateTemp("", "dst")
		require.NoError(t, err)
		defer os.Remove(dstFile.Name())

		content := []byte("Hello, World!")
		_, err = srcFile.Write(content)
		require.NoError(t, err)

		err = Copy(srcFile.Name(), dstFile.Name(), 0, 0)
		require.NoError(t, err)

		result, err := os.ReadFile(dstFile.Name())
		require.NoError(t, err)
		require.Equal(t, content, result)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		srcFile, err := os.CreateTemp("", "src")
		require.NoError(t, err)
		defer os.Remove(srcFile.Name())

		_, err = srcFile.Write([]byte("Hello"))
		require.NoError(t, err)

		dstFile, err := os.CreateTemp("", "dst")
		require.NoError(t, err)
		defer os.Remove(dstFile.Name())

		err = Copy(srcFile.Name(), dstFile.Name(), 10, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
}
