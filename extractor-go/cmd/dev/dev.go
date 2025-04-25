/**
 * "Development" entry point
 * This has limited scope to be used while developing something instead of running the entire history
 */

package main

import (
	"fmt"

	"github.com/guilherme-gm/ro-vis/extractor/internal/luaExtractor"
)

func main() {
	fmt.Println("RO Vis extractor - DEV")

	var quests []luaExtractor.QuestV1
	luaExtractor.Decode("../patches/OngoingQuestInfoList_True.lub", "QuestInfoList", &quests)

	fmt.Println(quests[0])

	fmt.Println("Success")
}
