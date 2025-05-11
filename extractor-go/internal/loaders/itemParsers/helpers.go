package itemParsers

import "github.com/guilherme-gm/ro-vis/extractor/internal/domain"

type txtEntry interface {
	GetItemID() int32
}

func loadTxtSubTable[T txtEntry](basePath string, update *domain.Update, newDB map[int32]*domain.Item, fileName string, parser func(string) (map[int32]T, error), mapper func(*domain.Item, *T)) {
	change, err := update.GetChangeForFile(fileName)
	if err != nil {
		// @TODO: Check if the error was actually the not found one
		return
	}

	itemVal, err := parser(basePath + "/" + change.Patch + "/" + fileName)
	if err != nil {
		panic(err)
	}

	for _, item := range newDB {
		if val, ok := itemVal[item.ItemID]; ok {
			mapper(item, &val)
		} else {
			mapper(item, nil)
		}
	}
}
