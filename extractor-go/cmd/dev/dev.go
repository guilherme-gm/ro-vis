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

	// loaders.ExtractInitialPatchList()

	patches, err := repository.GetPatchRepository().ListPatches()
	if err != nil {
		panic(err)
	}

	latest, err := repository.GetLoaderControllerRepository().GetLatestPatch("quest")
	if err != nil {
		panic(err)
	}

	loader := loaders.NewQuestLoader()
	for _, patch := range *patches {
		if patch.Id <= latest {
			continue
		}

		fmt.Println("Extracting " + patch.Name + "...")
		loader.LoadPatch(patch)

		repository.GetLoaderControllerRepository().SetLatestPatch("quest", patch.Id, patch.Name)
	}

	fmt.Println("Success")
}
