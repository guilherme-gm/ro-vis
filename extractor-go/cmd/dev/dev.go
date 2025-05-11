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
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders"
)

func dbCheck() {
	migTool, err := database.NewMigrationTool()
	if err != nil {
		panic(err)
	}

	migTool.UpdateCheck()
}

type loader interface {
	LoadPatch(basePath string, update domain.Update)
}

func load(updates []domain.Update, loaderName string, loaderInstance loader) {
	latest, err := repository.GetLoaderControllerRepository().GetLatestUpdate(loaderName)
	if err != nil {
		panic(err)
	}

	for _, update := range updates {
		if update.Date.Compare(latest) <= 0 {
			continue
		}

		fmt.Println("Extracting " + update.Name() + "...")
		loaderInstance.LoadPatch("../patches/", update)

		repository.GetLoaderControllerRepository().SetLatestPatch(loaderName, update.Date)
	}
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

	// load(updates, "quest", loaders.NewQuestLoader())
	load(updates, "items", loaders.NewItemLoader())

	fmt.Println("Success")
}
