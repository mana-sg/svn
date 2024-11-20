package repository

import (
	"database/sql"
	"fmt"

	"github.com/mana-sg/vcs/internal/db"
	"github.com/mana-sg/vcs/pkg/types"
)

type Repo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CommitCount int    `json:"commitCount"`
	UserId      int    `json:"userId"`
	UserName    string `json:"userName"`
}

type Commit struct {
	ID             int    `json:"id"`
	Message        string `json:"message"`
	TimeStamp      string `json:"timeStamp"`
	RepoID         int    `json:"repoId"`
	ParentCommitID *int   `json:"parentCommitId"`
}

func GetAllRepositories(db db.DbHandler, userId string) ([]Repo, error) {
	var repos []Repo

	queryString := `
    SELECT 
        repo.id AS repoId, 
        repo.name AS repoName, 
        repo.userId AS userId, 
        users.name AS userName, 
        (SELECT COUNT(*) 
        FROM vcs.commit 
        WHERE repoId = repo.id) AS commitCount
    FROM vcs.repo AS repo
    JOIN vcs.users AS users ON repo.userId = users.id
    WHERE repo.userId != ?
    AND repo.access = 1;
  `
	rows, err := db.GetValue(queryString, userId)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var repo Repo
		if err := rows.Scan(&repo.ID, &repo.Name, &repo.UserId, &repo.UserName, &repo.CommitCount); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		repos = append(repos, repo)
	}

	return repos, nil
}

func GetAllRepositoriesForUser(db db.DbHandler, userId string) ([]Repo, error) {
	queryString := `
    SELECT repo.id AS repoId, repo.name AS repoName, repo.userId as userId, 
          (SELECT COUNT(*) 
            FROM vcs.commit 
            WHERE repoId = repo.id) AS commitCount
    FROM vcs.repo AS repo
    WHERE repo.userId = ?
  `
	rows, err := db.GetValue(queryString, userId)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var repos []Repo

	for rows.Next() {
		var repo Repo
		if err := rows.Scan(&repo.ID, &repo.Name, &repo.UserId, &repo.CommitCount); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		repos = append(repos, repo)
	}

	return repos, nil
}

func GetNumberOfrepositories(db db.DbHandler, userId string) (int, error) {
	var numberOfRepos int
	query := `
    SELECT COUNT(*) AS repo_count
    FROM vcs.repo
    WHERE userId = ?;
  `
	row, err := db.GetValue(query, userId)
	if err != nil {
		return 0, fmt.Errorf(err.Error())
	}

	row.Next()
	row.Scan(&numberOfRepos)
	row.Close()

	return numberOfRepos, nil
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

func GetAllFilesForCommit(db db.DbHandler, commitId int) ([]types.FileNode, error) {
	var files []types.FileNode

	// Step 1: Get the root tree hash for the given commitId
	queryRoot := "SELECT referencesTree FROM vcs.commit WHERE id=?"
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

func buildFileTree(db db.DbHandler, treeHash string) ([]types.FileNode, error) {
	var nodes []types.FileNode

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
		node := types.FileNode{
			Id:       int8(id),
			Type:     int8(typeInt),
			Name:     name,
			Children: nil, // To be filled recursively if the node is a directory
			Content:  "",  // To be filled if it's a file
		}

		// Step 4: If it's a directory (type == 2), recursively fetch its children
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

func GetLatestCommit(db db.DbHandler, userId string, repoId int) (int, error) {
	queryStringLatestCommit := `
		SELECT c.id 
		FROM vcs.commit c
		JOIN vcs.repo r ON c.repoId = r.id
		WHERE r.userId = ? AND r.id = ?
		ORDER BY c.timeStamp DESC
		LIMIT 1;
	`

	rows, err := db.GetValue(queryStringLatestCommit, userId, repoId)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var commitId int
	if rows.Next() {
		if err := rows.Scan(&commitId); err != nil {
			return 0, fmt.Errorf("failed to scan commit ID: %v", err)
		}
	} else {
		return 0, fmt.Errorf("no commits found for user %d in repo %d", userId, repoId)
	}

	return commitId, nil
}
