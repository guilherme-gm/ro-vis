package loaders_test

import (
	"testing"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders"
	"github.com/stretchr/testify/assert"
)

func TestItemLoader_GetRelevantFiles(t *testing.T) {
	// Create a new ItemLoader with default parsers
	loader := loaders.NewItemLoader(server.GetTestServer()) // We pass nil for server since we don't need it for this test

	// Get the list of relevant fileMatchers
	fileMatchers := loader.GetRelevantFiles()

	// We expect all the files from all item parsers
	expectedFiles := []string{
		// ItemV1Parser
		"data/num2itemdisplaynametable.txt",
		"data/num2itemresnametable.txt",
		"data/itemslotcounttable.txt",
		// ItemV2Parser - same as V1
		// ItemV3Parser - same as V1
		// ItemV4Parser - same as V1
		// ItemV5Parser - same as V1
		// ItemV6Parser - same as V1
		// ItemV7Parser
	}

	// Check that all expected files are in the result
	for _, expectedFile := range expectedFiles {
		found := false
		for _, fileMatcher := range fileMatchers {
			if fileMatcher.MatchString(expectedFile) {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected file '%s' to be matched by relevant files, but it was not", expectedFile)
	}

	// Check that there are no duplicate files
	fileMap := make(map[string]bool)
	for _, fileMatcher := range fileMatchers {
		assert.False(t, fileMap[fileMatcher.String()], "Duplicate file found in relevant files: %s", fileMatcher)
		fileMap[fileMatcher.String()] = true
	}
}
