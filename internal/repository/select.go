package repository

import (
	"github.com/mana-sg/vcs/pkg/models"
)

func SelectRepo(id uint, name string) {
	models.SetActiveRepo(id, name)
}
