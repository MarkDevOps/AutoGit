package api_test

import (
	"os"
	"testing"

	"github.com/MarkDevOps/AutoGit/cli/pkg/api"
)

// Writing test to return a value for function setHeader
func TestSetHeader(t *testing.T) {
	// Test case: GITHUB_TOKEN is set
	t.Run("GITHUB_TOKEN is set", func(t *testing.T) {
		expectedToken := "gh_testtoken"
		os.Setenv("GITHUB_TOKEN", expectedToken) // Set the environment variable
		defer os.Unsetenv("GITHUB_TOKEN")        // Ensure the environment variable is set

		result := api.SetHeader()
		if result != expectedToken {
			t.Errorf("Expected %s, got %s", expectedToken, result)
		}
	})

	// Test case: GITHUB_TOKEN is not set
	t.Run("GITHUB_TOKEN is not set", func(t *testing.T) {
		os.Unsetenv("GITHUB_TOKEN") // Ensure environment is not set

		result := api.SetHeader()
		expectedError := "Error!!"
		if result != expectedError {
			t.Errorf("Expected %s, got %s", expectedError, result)
		}
	})
}
