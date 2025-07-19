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

	// Get the list of relevant files
	files := loader.GetRelevantFiles()

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
		assert.Contains(t, files, expectedFile, "Expected file %s not found in relevant files", expectedFile)
	}

	// Check that there are no duplicate files
	fileMap := make(map[string]bool)
	for _, file := range files {
		assert.False(t, fileMap[file], "Duplicate file found in relevant files: %s", file)
		fileMap[file] = true
	}
}
