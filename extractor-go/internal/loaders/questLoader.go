package loaders

import (
	"fmt"
	"strings"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

func HasQuestFiles(patch domain.Patch) bool {
	for _, fname := range patch.Files {
		if fname == "System/OngoingQuestInfoList_True.lub" {
			return true
		}

		lowerName := strings.ToLower(fname)
		if lowerName == "system/ongoingquestinfolist_true.lub" {
			fmt.Println("FOUND on lower -- " + fname)
			return true
		}
	}

	return false
}

func ExtractQuests(patch domain.Patch) {
	patchPath := "../patches/" + patch.Name + "/"

	fmt.Println("> Decoding...")
	var fileQuests []rostructs.QuestV3
	decoders.DecodeLuaTable(patchPath+"System/OngoingQuestInfoList_True.lub", "QuestInfoList", &fileQuests)

	fmt.Println("> Fetching current list...")
	currentQuests, err := repository.GetQuestRepository().GetCurrentQuests()
	if err != nil {
		panic(err)
	}

	fmt.Println("> Diffing...")
	questMap := make(map[int32]*domain.Quest)
	deletedIds := make(map[int32]bool)
	for _, q := range *currentQuests {
		questMap[q.QuestID] = &q

		if !q.Deleted {
			deletedIds[q.QuestID] = true
		}
	}

	var newQuests []domain.Quest
	var updatedQuests []domain.Quest

	for _, roQuest := range fileQuests {
		delete(deletedIds, roQuest.QuestId)
		fileQuest := roQuest.ToDomain()
		existingQuest := questMap[roQuest.QuestId]
		if existingQuest == nil {
			newQuests = append(newQuests, fileQuest)
			continue
		}

		if !existingQuest.Equals(fileQuest) {
			fileQuest.PreviousHistoryID = existingQuest.HistoryID
			updatedQuests = append(updatedQuests, fileQuest)
		}
	}

	fmt.Printf("> Saving new records... (%d records to save)\n", len(newQuests))
	err = repository.GetQuestRepository().AddQuestsToHistory(patch.Name, &newQuests)
	if err != nil {
		panic(err)
	}

	fmt.Printf("> Updating records... (%d records to update)\n", len(updatedQuests))
	err = repository.GetQuestRepository().AddQuestsToHistory(patch.Name, &updatedQuests)
	if err != nil {
		panic(err)
	}

	fmt.Printf("> Deleting records... (%d records to delete)\n", len(deletedIds))
	for deletedId := range deletedIds {
		err := repository.GetQuestRepository().AddDeletedQuest(patch.Name, questMap[deletedId])
		if err != nil {
			panic(err)
		}
	}
}
