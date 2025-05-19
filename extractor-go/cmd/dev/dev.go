/**
 * "Development" entry point
 * This has limited scope to be used while developing something instead of running the entire history
 */

package main

import (
	"database/sql"
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
	LoadPatch(tx *sql.Tx, basePath string, update domain.Update)
}

func load(updates []domain.Update, loaderName string, loaderInstance loader) {
	latest, err := repository.GetLoaderControllerRepository().GetLatestUpdate(nil, loaderName)
	if err != nil {
		panic(err)
	}

	for _, update := range updates {
		if update.Date.Compare(latest) <= 0 {
			continue
		}

		tx, err := database.BeginTx()
		if err != nil {
			panic(err)
		}
		defer tx.Rollback()

		fmt.Println("Extracting " + update.Name() + "...")
		loaderInstance.LoadPatch(tx, "../patches/", update)

		repository.GetLoaderControllerRepository().SetLatestPatch(tx, loaderName, update.Date)

		if err := tx.Commit(); err != nil {
			panic(err)
		}
	}
}

func main() {
	fmt.Println("RO Vis extractor - DEV")
	conf.LoadExtractor()
	dbCheck()

	// loaders.ExtractInitialPatchList()

	updates, err := repository.GetPatchRepository().ListUpdates(nil)
	if err != nil {
		panic(err)
	}

	// load(updates, "quest", loaders.NewQuestLoader())
	load(updates, "items", loaders.NewItemLoader())

	fmt.Println("Success")
}
