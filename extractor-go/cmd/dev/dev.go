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
	"github.com/guilherme-gm/ro-vis/extractor/internal/luaExtractor"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
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

	var quests []rostructs.QuestV1
	luaExtractor.Decode("../patches/OngoingQuestInfoList_True.lub", "QuestInfoList", &quests)

	var model = quests[0].ToModel()
	err := repository.GetQuestRepository().AddQuestToHistory(&model)
	if err != nil {
		panic(err)
	}

	fmt.Println(quests[0])

	fmt.Println("Success")
}
