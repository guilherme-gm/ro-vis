package extractor

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/luaExtractor"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

func ExtractQuest() {
	var fileQuests []rostructs.QuestV1
	luaExtractor.Decode("../patches/OngoingQuestInfoList_True.lub", "QuestInfoList", &fileQuests)

	currentQuests, err := repository.GetQuestRepository().GetCurrentQuests()
	if err != nil {
		panic(err)
	}

	questMap := make(map[int32]*domain.Quest)
	deletedIds := make(map[int32]bool)
	for _, q := range *currentQuests {
		questMap[q.QuestID] = &q
		deletedIds[q.QuestID] = true
	}

	var newQuests []domain.Quest
	var updatedQuests []domain.Quest

	for _, roQuest := range fileQuests {
		delete(deletedIds, roQuest.QuestId)
		fileQuest := roQuest.ToDomain()
		if questMap[roQuest.QuestId] == nil {
			newQuests = append(newQuests, fileQuest)
			continue
		}

		if !questMap[roQuest.QuestId].Equals(fileQuest) {
			updatedQuests = append(updatedQuests, fileQuest)
		}
	}

	for _, newQuest := range newQuests {
		err := repository.GetQuestRepository().AddQuestToHistory("2024-01-01", nil, &newQuest)
		if err != nil {
			panic(err)
		}
	}

	for _, updatedQuest := range updatedQuests {
		oldQuest := questMap[updatedQuest.QuestID]
		err := repository.GetQuestRepository().AddQuestToHistory("2024-01-01", oldQuest, &updatedQuest)
		if err != nil {
			panic(err)
		}
	}

	for deletedId := range deletedIds {
		err := repository.GetQuestRepository().AddDeletedQuest("2024-01-01", questMap[deletedId])
		if err != nil {
			panic(err)
		}
	}
}
