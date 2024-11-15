package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mana-sg/vcs/internal/db"
	"github.com/mana-sg/vcs/pkg/models"
)

func CreateRepo(db db.DbHandler, name string, userId string) error {
	if strings.Compare(name, "") == 0 {
		return fmt.Errorf("name cannot be empty")
	}
	repos, err := GetAllRepositoriesForUser(db, userId)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		if strings.Compare(repo.Name, name) == 0 {
			return fmt.Errorf("cannot create repository with same name, please use another unique name")
		}
	}

	createRepoQuery := "INSERT INTO vcs.repo(name, timeCreation, userId) VALUES(?, ?, ?)"

	userIdnum, err := strconv.Atoi(userId)
	if err != nil {
		return fmt.Errorf("error converting userId to int")
	}
	res, err := db.SetValue(createRepoQuery, name, time.Now(), userIdnum)
	if err != nil {
		return fmt.Errorf("error inserting into repo")
	}
	insId, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	models.SetActiveRepo(uint(insId), name)

	return nil
}
