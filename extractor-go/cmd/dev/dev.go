/**
 * "Development" entry point
 * This has limited scope to be used while developing something instead of running the entire history
 */

package main

import (
	"fmt"

	"github.com/guilherme-gm/ro-vis/extractor/internal/conf"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders"
)

func dbCheck() {
	migTool, err := database.NewMigrationTool()
	if err != nil {
		panic(err)
	}

	migTool.UpdateCheck()
}

func main() {
	fmt.Println("RO Vis extractor - DEV")
	conf.Load()
	dbCheck()

	patches, err := repository.GetPatchRepository().ListPatches()
	if err != nil {
		panic(err)
	}

	for _, patch := range *patches {
		fmt.Println("Extracting " + patch.Name + "...")
		if loaders.HasQuestFiles(patch) {
			loaders.ExtractQuests(patch)
			fmt.Println("Done")
		} else {
			fmt.Println("Skipped")
		}
	}
	// extractor.ExtractQuest()
	// extractor.ExtractInitialPatchList()

	fmt.Println("Success")
}
