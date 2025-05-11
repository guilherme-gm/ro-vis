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

	updates, err := repository.GetPatchRepository().ListUpdates()
	if err != nil {
		panic(err)
	}

	latest, err := repository.GetLoaderControllerRepository().GetLatestUpdate("quest")
	if err != nil {
		panic(err)
	}

	loader := loaders.NewQuestLoader()
	for _, update := range updates {
		if update.Date.Compare(latest) <= 0 {
			continue
		}

		fmt.Println("Extracting " + update.Name() + "...")
		loader.LoadPatch("../patches/", update)

		repository.GetLoaderControllerRepository().SetLatestPatch("quest", update.Date)
	}

	fmt.Println("Success")
}
