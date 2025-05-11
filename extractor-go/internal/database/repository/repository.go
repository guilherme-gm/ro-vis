package repository

type repositories struct {
	ItemRepository             *ItemRepository
	LoaderControllerRepository *LoaderControllerRepository
	patchRepository            *PatchRepository
	questRepository            *QuestRepository
}

var repositoriesCache repositories = repositories{}
