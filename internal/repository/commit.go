package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/mana-sg/vcs/internal/db"
	"github.com/mana-sg/vcs/internal/utils"
	"github.com/mana-sg/vcs/pkg/models"
	"github.com/mana-sg/vcs/pkg/types"
)

func CreateCommit(db db.DbHandler, message string, repoId int) (int, error) {
	//Get the active user name
	_, userName := models.GetActiveUser()

	//Check if the username is not initialised
	if strings.Compare(userName, "") == 0 {
		return 0, fmt.Errorf("user not selected")
	}

	//creating the query for committing into a repository
	createCommitQuery := "INSERT INTO vcs.commit(message, timeStamp, repoId) VALUES(?, ?, ?)"

	//executing the query and getting the error message
	insertion, err := db.SetValue(createCommitQuery, message, time.Now(), repoId)
	//Check for any errors in the process
	fmt.Println(err)
	if err != nil {
		return 0, fmt.Errorf("error adding commit: %v", err)
	}
	commitId, err := insertion.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting commit id: %v", err)
	}

	//if there is no error return nil
	return int(commitId), nil
}

func AddTreeEntry(db db.DbHandler, node types.FileNode, parentId string) error {
	emptyHash := sha256.Sum256([]byte(""))

	if strings.Compare(node.Name, "") == 0 {
		return fmt.Errorf("name of file cannot be empty")
	}

	// Do not make an entry for an empty file
	if node.Type == 1 && strings.Compare(node.Content, hex.EncodeToString(emptyHash[:])) == 0 {
		return nil
	}

	// Hash the file or directory contents for a unique ID
	var hash []byte
	if node.Type == 1 {
		// File: hash the content
		hash, _ = utils.Hash([]byte(node.Content))
	} else {
		// Directory: recursively hash children
		hash, _ = utils.HashDirectoryContents(node.Children)
	}

	// Check if this hash already exists in the relevant table
	exists, err := hashExists(db, string(hash), node.Type)
	if err != nil {
		return fmt.Errorf("error checking hash existence: %v", err)
	}
	if exists {
		// If the hash already exists, just link to it in the tree_entry table
		return linkExistingTreeEntry(db, node, parentId, string(hash))
	}

	// Create a new entry in the tree_entry table
	queryString := "INSERT INTO vcs.tree_entry(type, name, parentTreeId) VALUES(?, ?, ?)"
	treeEntryInsertion, err := db.SetValue(queryString, node.Type, node.Name, parentId)
	if err != nil {
		return fmt.Errorf("error inserting tree entry: %v", err)
	}
	treeEntryId, err := treeEntryInsertion.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting tree entry id")
	}

	if node.Type == 1 {
		// If it's a file, insert the content into blobContent
		queryStringEntry := "INSERT INTO vcs.blobContent(hash, content) VALUES(?, ?)"
		_, err = db.SetValue(queryStringEntry, string(hash), node.Content)
		if err != nil {
			return fmt.Errorf("error inserting blob entry: %v", err)
		}
		linkBlobQuery := "UPDATE vcs.tree_entry SET childBlobId=? WHERE id = ?"
		_, err := db.SetValue(linkBlobQuery, hash, treeEntryId)
		if err != nil {
			return fmt.Errorf("error updating child blob id")
		}
		fmt.Println("i am after update")
	} else if node.Type == 2 {
		// If it's a directory, add its subtree
		err := AddTree(db, string(hash), node.Children, 0)
		if err != nil {
			return err
		}
		linkTreeQuery := "UPDATE vcs.tree_entry SET childTreeId=? WHERE id = ?"
		_, err = db.SetValue(linkTreeQuery, hash, treeEntryId)
		if err != nil {
			return fmt.Errorf("error updating child tree id")
		}
	}

	return nil
}

// Recursive function to add a tree (directory)
func AddTree(db db.DbHandler, dirHash string, children []types.FileNode, commitId int) error {
	// Do not make an entry for an empty directory
	emptyHash := sha256.Sum256([]byte(""))
	if strings.Compare(string(dirHash), string(hex.EncodeToString(emptyHash[:]))) == 0 {
		return nil
	}

	// Check if this directory hash already exists in the tree table
	exists, err := hashExists(db, dirHash, 2)
	if err != nil {
		return fmt.Errorf("error checking tree hash existence: %v", err)
	}
	if !exists {
		// Insert the directory hash into the tree table
		queryString := "INSERT INTO vcs.tree(hash) VALUES(?)"
		_, err = db.SetValue(queryString, dirHash)
		fmt.Println(err)
		if err != nil {
			return fmt.Errorf("error inserting tree entry: %v", err)
		}
	}

	fmt.Println("here outside for commitId", commitId)
	if commitId != 0 {
		fmt.Println("here for commitId", commitId)
		queryStringUpdate := "UPDATE vcs.commit SET referencesTree=? WHERE id=?"
		fmt.Println(err)
		_, err := db.SetValue(queryStringUpdate, dirHash, commitId)
		if err != nil {
			fmt.Errorf("error in adding refernce to tree")
		}
	}
	// Insert each child in the directory
	for _, child := range children {
		err := AddTreeEntry(db, child, dirHash)
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function to check if a hash already exists in the relevant table
func hashExists(db db.DbHandler, hash string, entryType int8) (bool, error) {
	var table string
	if entryType == 1 {
		table = "vcs.blobContent"
	} else {
		table = "vcs.tree"
	}

	queryString := fmt.Sprintf("SELECT hash FROM %s WHERE hash = ?", table)
	rows, err := db.GetValue(queryString, hash)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Check if any row is returned (indicating the hash exists)
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

// Helper function to link an existing tree entry without creating a new one
func linkExistingTreeEntry(db db.DbHandler, node types.FileNode, hash string, childHash string) error {
	fmt.Println("I am here for node: ", node.Name)
	queryString := "INSERT INTO vcs.tree_entry(type, name, parentTreeId) VALUES(?, ?, ?)"

	treeEntryInsertion, err := db.SetValue(queryString, node.Type, node.Name, hash)
	if err != nil {
		return fmt.Errorf("error inserting tree entry: %v", err)
	}
	fmt.Println("got tree_insertion")
	treeEntryId, err := treeEntryInsertion.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting tree entry id: %v", err)
	}
	fmt.Println("here with tree_entry id", treeEntryId)

	if node.Type == 1 {
		// If it's a file, insert the content into blobContent
		linkBlobQuery := "UPDATE vcs.tree_entry SET childBlobId=? WHERE id = ?"
		_, err := db.SetValue(linkBlobQuery, childHash, treeEntryId)
		fmt.Println(err)
		if err != nil {
			return fmt.Errorf("error updating child blob id")
		}
		fmt.Println("ok no error")
	} else if node.Type == 2 {
		// If it's a directory, add its subtree
		err := AddTree(db, string(childHash), node.Children, 0)
		if err != nil {
			return err
		}
		linkTreeQuery := "UPDATE vcs.tree_entry SET childTreeId=? WHERE id = ?"
		_, err = db.SetValue(linkTreeQuery, childHash, treeEntryId)
		if err != nil {
			return fmt.Errorf("error updating child tree id")
		}
	}
	fmt.Println("returning from linking")
	return nil
}
