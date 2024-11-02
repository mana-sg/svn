package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/mana-sg/vcs/internal/db"
	"github.com/mana-sg/vcs/pkg/models"
	"github.com/mana-sg/vcs/pkg/types/entry"
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

func AddTreeEntry(db db.DbHandler, entry entry.EntryArray) error {
	emptyhash := sha256.Sum256([]byte(""))
	if strings.Compare(entry.EntryName, "") == 0 {
		return fmt.Errorf("name of file cannot be the empty")
	}

	//Do not make an entry for an empty file
	if strings.Compare(entry.EntryHash, hex.EncodeToString(emptyhash[:])) == 0 {
		return nil
	}

	//creating query string for treee entry which consists of id(hash),the type of entry(0 file, 1 dir) and name of the file/dir
	queryString := "INSERT INTO vcs.tree_entry(id, type, name) VALUES(? ,? ,?)"
	//setting value into database
	_, err := db.SetValue(queryString, entry.EntryHash, entry.EntryType, entry.EntryName)

	//error handling
	if err != nil {
		return fmt.Errorf("error inserting tree entry: %v", err)
	}

	//checking if it is a file or directory
	if entry.EntryType == 0 {
		//creating query string
		queryStringEntry := "INSERT INTO vcs.blobContent(hash, content) VALUES(? ,?)"

		//inserting value into blobContent table
		_, err := db.SetValue(queryStringEntry, entry.EntryHash, entry.EntryContent)
		if err != nil {
			return fmt.Errorf("error inserting blob entry: %v", err)
		}
	} else if entry.EntryType == 1 {
		//if it is a file call the AddTree funciton which handles with adding tree
		err := AddTree(db, entry.EntryHash, entry.EntriesUnder)
		if err != nil {
			return err
		}
	}

	//return nil if no errorrs
	return nil
}

// TODO: adding pointers to already existing hashes if it already exists, to reduce redundancy
func AddTree(db db.DbHandler, totalHash string, entriesUnder []entry.EntryArray) error {
	emptyhash := sha256.Sum256([]byte(""))
	//Do not make an entry for an empty file
	if strings.Compare(totalHash, hex.EncodeToString(emptyhash[:])) == 0 {
		return nil
	}

	//creating queryString
	queryString := "INSERT INTO vcs.tree(hash) VALUES(?)"

	//inserting it into the tree table
	_, err := db.SetValue(queryString, totalHash)

	//error handling
	if err != nil {
		return fmt.Errorf("error inserting tree entry: %v", err)
	}

	//going through all the tree entries under the directory
	for _, entry := range entriesUnder {
		//calling the AddTreeEntry function for each entry under the dir
		err := AddTreeEntry(db, entry)
		//error handling
		if err != nil {
			return err
		}
	}

	//return nil if no errors
	return nil
}
