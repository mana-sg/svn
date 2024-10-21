package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/mana-sg/vcs/internal/db"
	"github.com/mana-sg/vcs/pkg/models"
)

func CreateCommit(db db.DbHandler, message string) error {
	_, userName := models.GetActiveUser()
	repoId, _ := models.GetActiveUser()
	if strings.Compare(userName, "") == 0 {
		return fmt.Errorf("user not selected")
	}

	createCommitQuery := "INSERT INTO vcs.commit(message, timeStamp, repoId) VALUES(?, ?, ?)"

	_, err := db.SetValue(createCommitQuery, message, time.Now().UnixNano(), repoId)
	if err != nil {
		return fmt.Errorf("error adding commit: %v", err)
	}

	return nil
}
