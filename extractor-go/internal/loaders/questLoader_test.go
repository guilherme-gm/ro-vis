package loaders_test

import (
	"testing"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders"
	"github.com/stretchr/testify/assert"
)

func TestQuestLoader_GetRelevantFiles(t *testing.T) {
	// Create a new QuestLoader with default parsers
	loader := loaders.NewQuestLoader(server.GetTestServer()) // We pass nil for server since we don't need it for this test

	// Get the list of relevant files
	files := loader.GetRelevantFiles()

	// We expect all the files from all quest parsers
	expectedFiles := []string{
		// QuestV1Parser
		"data/questid2display.txt",
		// QuestV2Parser - same as V1
		// QuestV3Parser
		"System/OngoingQuestInfoList_True.lub",
		// QuestV4Parser - same as V3
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
