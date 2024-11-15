package repository

import (
	"database/sql"
	"fmt"

	"github.com/mana-sg/vcs/internal/db"
)

type Repo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	TimeCreation string `json:"time_creation"`
}

type Commit struct {
	ID             int    `json:"id"`
	Message        string `json:"message"`
	TimeStamp      string `json:"timeStamp"`
	RepoID         int    `json:"repoId"`
	ParentCommitID *int   `json:"parentCommitId"`
}

type FileNode struct {
	Id       int8       `json:"id"`
	Name     string     `json:"name"`
	Type     int8       `json:"type"`
	Children []FileNode `json:"children"`
	Content  string     `json:"content"`
}

func GetAllRepositoriesForUser(db db.DbHandler, userId string) ([]Repo, error) {
	queryString := "SELECT id, name, timeCreation FROM vcs.repo where userId = ?"
	rows, err := db.GetValue(queryString, userId)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var repos []Repo

	for rows.Next() {
		var repo Repo
		if err := rows.Scan(&repo.ID, &repo.Name, &repo.TimeCreation); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		repos = append(repos, repo)
	}

	return repos, nil
}

func GetAllCommitsForRepo(db db.DbHandler, repoId string) ([]Commit, error) {
	var commits []Commit
	var lastCommit Commit
	queryGetCurrentCommit := `
    SELECT id, message, timeStamp, repoId, parentCommitId
    FROM vcs.commit
    WHERE repoId = ?
    ORDER BY timeStamp DESC
    LIMIT 1;
  `
	row, err := db.GetValue(queryGetCurrentCommit, repoId)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	row.Next()
	row.Scan(&lastCommit.ID, &lastCommit.Message, &lastCommit.TimeStamp, &lastCommit.RepoID, &lastCommit.ParentCommitID)
	row.Close()

	for {
		commits = append(commits, lastCommit)

		if lastCommit.ParentCommitID == nil {
			break
		}
		queryPrevCommit := `
      SELECT id, message, timeStamp, repoId, parentCommitId
      FROM vcs.commit 
      WHERE id=?
    `
		nextRow, err := db.GetValue(queryPrevCommit, lastCommit.ParentCommitID)
		defer nextRow.Close()
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		nextRow.Next()
		nextRow.Scan(&lastCommit.ID, &lastCommit.Message, &lastCommit.TimeStamp, &lastCommit.RepoID, &lastCommit.ParentCommitID)
	}

	return commits, nil
}

func GetAllFilesForCommit(db db.DbHandler, commitId int) ([]FileNode, error) {
	var files []FileNode

	// Step 1: Get the root tree hash for the given commitId
	queryRoot := "SELECT hash FROM vcs.tree WHERE pointsToCommit=?"
	var rootHash string
	row, err := db.GetValue(queryRoot, commitId)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	row.Next()
	row.Scan(&rootHash)

	// Step 2: Build the file tree from the root tree hash
	files, err = buildFileTree(db, rootHash)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func buildFileTree(db db.DbHandler, treeHash string) ([]FileNode, error) {
	var nodes []FileNode

	// Step 3: Get the tree entries for this tree
	query := "SELECT id, name, type, childBlobId, childTreeId FROM vcs.tree_entry WHERE parentTreeId=?"
	rows, err := db.GetValue(query, treeHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var typeInt int
		var childBlobId, childTreeId sql.NullString
		var id int8

		// Read the tree entry details
		err := rows.Scan(&id, &name, &typeInt, &childBlobId, &childTreeId)
		if err != nil {
			return nil, err
		}

		// Create the FileNode
		node := FileNode{
			Id:       int8(id),
			Type:     int8(typeInt),
			Name:     name,
			Children: nil, // To be filled recursively if the node is a directory
			Content:  "",  // To be filled if it's a file
		}

		// Step 4: If it's a directory (type == 1), recursively fetch its children
		if typeInt == 2 { // Directory
			if childTreeId.String != "" {
				// Recursive call for directory
				node.Children, err = buildFileTree(db, childTreeId.String)
				if err != nil {
					return nil, err
				}
			}
		} else if typeInt == 1 { // File
			if childBlobId.String != "" {
				// Get the file content from blobContent
				contentQuery := "SELECT content FROM vcs.blobContent WHERE hash=?"
				var content []byte
				row, err := db.GetValue(contentQuery, childBlobId.String)
				if err != nil {
					return nil, err
				}
				defer row.Close()
				row.Next()
				row.Scan(&content)

				node.Content = string(content)
			}
		}

		// Add the node to the list of nodes
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func GetLatestCommit(userId int, repoId int) {

}
