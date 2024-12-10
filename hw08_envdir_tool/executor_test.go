package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("run command with environment variables", func(t *testing.T) {
		env := Environment{
			"FOO": {Value: "123", NeedRemove: false},
			"BAR": {Value: "value", NeedRemove: false},
		}

		// Create a temporary Go file to print environment variables
		tempDir, err := os.MkdirTemp("", "envtest")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		goFilePath := filepath.Join(tempDir, "printenv.go")
		goFileContent := `
package main

import (
    "fmt"
    "os"
)

func main() {
    for _, e := range os.Environ() {
        fmt.Println(e)
    }
}
`
		err = os.WriteFile(goFilePath, []byte(goFileContent), 0o644)
		require.NoError(t, err)

		// Use `go run` to execute the temporary Go file
		cmd := []string{"go", "run", goFilePath}
		exitCode := RunCmd(cmd, env)
		require.Equal(t, 0, exitCode)
	})
}
