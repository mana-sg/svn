package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/mana-sg/vcs/internal/db"
	"github.com/mana-sg/vcs/pkg/models"
)

func CreateCommit(db db.DbHandler, message string) error {
	//Get the active user name
	_, userName := models.GetActiveUser()

	//Get the active repo id
	repoId, _ := models.GetActiveUser()

	//Check if the username is not initialised
	if strings.Compare(userName, "") == 0 {
		return fmt.Errorf("user not selected")
	}

	//creating the query for committing into a repository
	createCommitQuery := "INSERT INTO vcs.commit(message, timeStamp, repoId) VALUES(?, ?, ?)"

	//executing the query and getting the error message
	_, err := db.SetValue(createCommitQuery, message, time.Now().UnixNano(), repoId)
	//Check for any errors in the process
	if err != nil {
		return fmt.Errorf("error adding commit: %v", err)
	}

	//if there is no error return nil
	return nil
}

func AddTreeEntry(db db.DbHandler, hash string, name string, typeEntry uint8) error {
	emptyhash := sha256.Sum256([]byte(""))
	if strings.Compare(name, "") == 0 {
		return fmt.Errorf("name of file cannot be the empty")
	}

	//Do not make an entry for an empty file
	if strings.Compare(hash, hex.EncodeToString(emptyhash[:])) == 0 {
		return nil
	}

	return nil
}
