package repository

import (
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
	defer row.Close()
	row.Scan(&lastCommit.ID, &lastCommit.Message, &lastCommit.TimeStamp, &lastCommit.RepoID, &lastCommit.ParentCommitID)

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
		row, err := db.GetValue(queryPrevCommit, lastCommit.ParentCommitID)
		defer row.Close()
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		row.Scan(&lastCommit.ID, &lastCommit.Message, &lastCommit.TimeStamp, &lastCommit.RepoID, &lastCommit.ParentCommitID)
	}

	return commits, nil
}

func GetAllFilesForCommit(db db.DbHandler, commitId string) ([]File, error) {
	var files []File

	return files, nil
}
