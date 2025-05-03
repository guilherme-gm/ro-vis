package repository

type repositories struct {
	patchRepository *PatchRepository
	questRepository *QuestRepository
}

var repositoriesCache repositories = repositories{}
