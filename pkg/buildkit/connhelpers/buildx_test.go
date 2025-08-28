package connhelpers

import (
	"net/url"
	"os/exec"
	"testing"

	"github.com/moby/buildkit/client/connhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildx(t *testing.T) {
	_, err := connhelper.GetConnectionHelper("buildx://")
	assert.NoError(t, err)

	_, err = connhelper.GetConnectionHelper("buildx://foobar")
	assert.NoError(t, err)

	_, err = connhelper.GetConnectionHelper("buildx://foorbar/something")
	assert.Error(t, err)
}

func TestBuildxWithPath(t *testing.T) {
	// Test that Buildx function properly rejects paths
	u, err := url.Parse("buildx://builder/path")
	require.NoError(t, err)

	_, err = Buildx(u)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not support path elements")
}

func TestBuildxWithoutPath(t *testing.T) {
	// Test that Buildx function works without path
	u, err := url.Parse("buildx://builder")
	require.NoError(t, err)

	helper, err := Buildx(u)
	assert.NoError(t, err)
	assert.NotNil(t, helper)
	assert.NotNil(t, helper.ContextDialer)
}

func TestBuildxHelperCreation(t *testing.T) {
	// Test that the connection helper is properly created without executing it
	tests := []struct {
		name       string
		builder    string
		expectFunc bool
	}{
		{
			name:       "with named builder",
			builder:    "test-builder",
			expectFunc: true,
		},
		{
			name:       "with empty builder",
			builder:    "",
			expectFunc: true,
		},
		{
			name:       "with complex builder name",
			builder:    "my-complex-builder-name",
			expectFunc: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse buildx URL
			u, err := url.Parse("buildx://" + tt.builder)
			require.NoError(t, err)

			// Get connection helper - should not execute anything
			helper, err := Buildx(u)
			require.NoError(t, err)
			require.NotNil(t, helper)

			if tt.expectFunc {
				require.NotNil(t, helper.ContextDialer, "ContextDialer should not be nil")
			}
		})
	}
}

// Helper function to check if docker is available
func isDockerAvailable() bool {
	cmd := exec.Command("docker", "--version")
	return cmd.Run() == nil
}
