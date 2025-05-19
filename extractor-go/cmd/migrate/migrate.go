/**
 * "Migrate" entry point
 * Applies DB migrations
 */

package main

import (
	"fmt"

	"github.com/guilherme-gm/ro-vis/extractor/internal/conf"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
)

func main() {
	fmt.Println("RO Vis extractor - Migrate")
	conf.LoadExtractor()

	migTool, err := database.NewMigrationTool()
	if err != nil {
		panic(err)
	}

	err = migTool.Up()
	if err != nil {
		fmt.Println("Failed to update")
		panic(err)
	}

	fmt.Println("Success")
}
