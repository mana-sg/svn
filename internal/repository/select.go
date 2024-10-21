package repository

import (
	"fmt"

	"github.com/mana-sg/vcs/internal/db"
	"github.com/mana-sg/vcs/pkg/models"
)

func GetAllRepos(db db.DbHandler) ([]struct {
	id   uint
	name string
}, error) {
	var repos []struct {
		id   uint
		name string
	}

	userId, _ := models.GetActiveUser()

	getAllReposQuery := "SELECT id, name FROM vcs.repo WHERE userId=?"
	allRepos, err := db.GetValue(getAllReposQuery, userId)
	if err != nil {
		return []struct {
			id   uint
			name string
		}{}, err
	}
	defer allRepos.Close()

	for allRepos.Next() {
		var repo struct {
			id   uint
			name string
		}
		err := allRepos.Scan(&repo.id, &repo.name)
		if err != nil {
			return []struct {
				id   uint
				name string
			}{}, fmt.Errorf("error retrieving repo name: %v", err)
		}

		repos = append(repos, repo)
	}
	return repos, nil
}

func SelectRepo(id uint, name string) {
	models.SetActiveRepo(id, name)
}
