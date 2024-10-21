//This model is to store info on which repository is active while the program is running

package models

type repo struct {
	id   uint
	name string
}

var activeRepo repo

func SetActiveRepo(id uint, name string) {
	activeRepo.id = id
	activeRepo.name = name
}

func GetActiveRepo() (uint, string) {
	return activeRepo.id, activeRepo.name
}
