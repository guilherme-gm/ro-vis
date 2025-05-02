package repository

type repositories struct {
	questRepository *QuestRepository
}

var repositoriesCache repositories = repositories{}
