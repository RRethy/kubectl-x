package build_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/RRethy/utils/kustomizelite/pkg/kustomize"
	"github.com/RRethy/utils/kustomizelite/pkg/testutil"
)

func TestKustomizerIntegrationWithFixtures(t *testing.T) {
	k, err := kustomize.NewKustomize(nil)
	require.NoError(t, err)

	t.Run("valid fixtures", func(t *testing.T) {
		fixtures := []struct {
			name     string
			path     string
			validate func(t *testing.T, content []map[string]any)
		}{
			{
				name: "basic kustomization",
				path: testutil.ValidFixturePath("basic"),
				validate: func(t *testing.T, content []map[string]any) {
					t.Helper()
					require.Equal(t, 2, len(content)) // 2 resources only (kustomization excluded)
					// Check we have deployment and service
					kinds := make(map[string]bool)
					for _, doc := range content {
						kinds[doc["kind"].(string)] = true
					}
					assert.True(t, kinds["Deployment"])
					assert.True(t, kinds["Service"])
				},
			},
			{
				name: "with namespace and labels",
				path: testutil.ValidFixturePath("with-namespace"),
				validate: func(t *testing.T, content []map[string]any) {
					t.Helper()
					require.Equal(t, 1, len(content)) // 1 configmap only (kustomization excluded)
					// Check we have the configmap
					assert.Equal(t, "ConfigMap", content[0]["kind"])
					metadata := content[0]["metadata"].(map[string]any)
					assert.Equal(t, "namespace-config", metadata["name"])
					// Verify namespace was applied
					assert.Equal(t, "production", metadata["namespace"])
				},
			},
			{
				name: "with patches",
				path: testutil.ValidFixturePath("with-patches"),
				validate: func(t *testing.T, content []map[string]any) {
					t.Helper()
					require.Equal(t, 2, len(content)) // 2 resources only (kustomization excluded)
					// Check we have deployment and service
					kinds := make(map[string]bool)
					for _, doc := range content {
						kinds[doc["kind"].(string)] = true
					}
					assert.True(t, kinds["Deployment"])
					assert.True(t, kinds["Service"])
				},
			},
			{
				name: "with components",
				path: testutil.ValidFixturePath("with-components"),
				validate: func(t *testing.T, content []map[string]any) {
					t.Helper()
					require.Equal(t, 4, len(content)) // 2 from resources + 2 from components
					// Check we have all expected resource types
					kinds := make(map[string]bool)
					for _, doc := range content {
						kinds[doc["kind"].(string)] = true
					}
					assert.True(t, kinds["Deployment"])
					assert.True(t, kinds["Service"])
					assert.True(t, kinds["ConfigMap"]) // From monitoring component
					assert.True(t, kinds["Secret"])    // From security component
				},
			},
			{
				name: "empty kustomization",
				path: testutil.ValidFixturePath("empty"),
				validate: func(t *testing.T, content []map[string]any) {
					t.Helper()
					require.Len(t, content, 0) // Empty since kustomization has no resources
				},
			},
		}

		for _, tt := range fixtures {
			t.Run(tt.name, func(t *testing.T) {
				content, err := k.Kustomize(tt.path, nil)
				require.NoError(t, err)
				tt.validate(t, content)
			})
		}
	})

	t.Run("invalid fixtures", func(t *testing.T) {
		fixtures := []struct {
			name    string
			path    string
			wantErr string
		}{
			{
				name:    "malformed yaml",
				path:    testutil.InvalidFixturePath("malformed-yaml"),
				wantErr: "parsing Kustomization YAML",
			},
			{
				name:    "invalid syntax",
				path:    testutil.InvalidFixturePath("invalid-syntax"),
				wantErr: "parsing Kustomization YAML",
			},
			{
				name:    "tab indentation",
				path:    testutil.InvalidFixturePath("tab-indentation"),
				wantErr: "parsing Kustomization YAML",
			},
			{
				name:    "duplicate keys",
				path:    testutil.InvalidFixturePath("duplicate-keys"),
				wantErr: "parsing Kustomization YAML",
			},
		}

		for _, tt := range fixtures {
			t.Run(tt.name, func(t *testing.T) {
				_, err := k.Kustomize(tt.path, nil)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			})
		}
	})

	t.Run("recursive resource loading", func(t *testing.T) {
		fixtures := []struct {
			name     string
			path     string
			validate func(t *testing.T, content []map[string]any)
		}{
			{
				name: "with resources",
				path: testutil.ValidFixturePath("with-resources"),
				validate: func(t *testing.T, content []map[string]any) {
					t.Helper()
					// Should have loaded resources only (kustomization excluded)
					assert.Equal(t, 2, len(content)) // deployment + service

					// Check we have all expected kinds
					kinds := make(map[string]bool)
					for _, doc := range content {
						kinds[doc["kind"].(string)] = true
					}
					assert.False(t, kinds["Kustomization"]) // Should not include kustomization
					assert.True(t, kinds["Deployment"])
					assert.True(t, kinds["Service"])
				},
			},
			{
				name: "with bases (bases no longer processed)",
				path: testutil.ValidFixturePath("with-bases"),
				validate: func(t *testing.T, content []map[string]any) {
					t.Helper()
					// Should have loaded ingress only (kustomization excluded, bases ignored)
					assert.Equal(t, 1, len(content)) // ingress only

					// Count document types
					kustomizationCount := 0
					ingressCount := 0
					for _, doc := range content {
						switch doc["kind"] {
						case "Kustomization":
							kustomizationCount++
						case "Ingress":
							ingressCount++
						}
					}
					assert.Equal(t, 0, kustomizationCount, "Should have 0 Kustomization documents (excluded from results)")
					assert.Equal(t, 1, ingressCount, "Should have 1 Ingress document")
				},
			},
		}

		for _, tt := range fixtures {
			t.Run(tt.name, func(t *testing.T) {
				content, err := k.Kustomize(tt.path, nil)
				require.NoError(t, err)
				tt.validate(t, content)
			})
		}
	})
}

func TestDuplicateKeysHandling(t *testing.T) {
	k, err := kustomize.NewKustomize(nil)
	require.NoError(t, err)

	// The duplicate-keys fixture has duplicate namespace and resources keys
	// YAML parsers typically use the last value for duplicate keys
	content, err := k.Kustomize(testutil.InvalidFixturePath("duplicate-keys"), nil)

	// Depending on the YAML parser's behavior, this might succeed or fail
	// Let's test what actually happens
	if err != nil {
		assert.Contains(t, err.Error(), "parsing Kustomization YAML")
	} else {
		// If it succeeds, should have only the service.yaml resource (kustomization excluded)
		// The duplicate-keys fixture has resources: ["service.yaml"] as the last value
		require.Len(t, content, 1)
		assert.Equal(t, "Service", content[0]["kind"])
	}
}

func TestKustomizerWithDirectoryPaths(t *testing.T) {
	k, err := kustomize.NewKustomize(nil)
	require.NoError(t, err)

	t.Run("directory paths work same as file paths", func(t *testing.T) {
		fixtures := []string{
			"basic",
			"with-namespace",
			"with-patches",
			"empty",
		}

		for _, fixture := range fixtures {
			t.Run(fixture, func(t *testing.T) {
				// Test with file path
				filePath := testutil.ValidFixturePath(fixture)
				fileContent, fileErr := k.Kustomize(filePath, nil)

				// Test with directory path
				dirPath := filepath.Dir(filePath)
				dirContent, dirErr := k.Kustomize(dirPath, nil)

				// Both should succeed
				require.NoError(t, fileErr)
				require.NoError(t, dirErr)

				// Both should return the same content
				assert.Equal(t, fileContent, dirContent)
			})
		}
	})

	t.Run("directory without kustomization.yaml fails", func(t *testing.T) {
		tmpDir := t.TempDir()
		_, err := k.Kustomize(tmpDir, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "stat'ing path")
	})
}
