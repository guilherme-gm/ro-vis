package repository

type repositories struct {
	LoaderControllerRepository *LoaderControllerRepository
	patchRepository            *PatchRepository
	questRepository            *QuestRepository
}

var repositoriesCache repositories = repositories{}
