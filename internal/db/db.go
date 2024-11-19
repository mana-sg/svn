package db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type DbHandler struct {
	Db *sql.DB
}

func (d *DbHandler) ConfigDB(user, password, host, dbname string) error {
	cfg := mysql.Config{
		User:            user,
		Passwd:          password,
		Net:             "tcp",
		Addr:            host,
		DBName:          dbname,
		MultiStatements: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}
	d.Db = db

	return nil
}

func (d *DbHandler) Close() error {
	err := d.Db.Close()
	if err != nil {
		return fmt.Errorf("error closing database: %v", err)
	}
	return nil
}

func (d *DbHandler) SetValue(query string, args ...interface{}) (sql.Result, error) {
	res, err := d.Db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error inserting values: %v", err)
	}
	return res, nil
}

func (d *DbHandler) GetValue(query string, args ...interface{}) (*sql.Rows, error) {
	res, err := d.Db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting the value: %v", err)
	}
	return res, nil
}

func (d *DbHandler) PrepareDb() error {
	createDbQuery := "CREATE DATABASE IF NOT EXISTS vcs"

	_, err := d.Db.Exec(createDbQuery)

	if err != nil {
		return fmt.Errorf("error creating database: %s", err)
	}

	_, err = d.Db.Exec("USE vcs")
	if err != nil {
		return fmt.Errorf("error selecting database: %s", err)
	}

	createTablesQuery := []string{`
    CREATE TABLE IF NOT EXISTS vcs.users(
      id INT AUTO_INCREMENT PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      email VARCHAR(255) NOT NULL UNIQUE,
      password VARCHAR(255) NOT NULL
    );`,
		`CREATE TABLE IF NOT EXISTS vcs.repo(
      id INT AUTO_INCREMENT PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      timeCreation DATETIME NOT NULL, 
      userId INT NOT NULL,
      FOREIGN KEY (userId) REFERENCES vcs.users(id)
    );`,
		`CREATE TABLE IF NOT EXISTS vcs.tree(
      hash VARCHAR(64) NOT NULL PRIMARY KEY
    );`,
		`CREATE TABLE IF NOT EXISTS vcs.commit(
      id INT AUTO_INCREMENT PRIMARY KEY,
      message VARCHAR(255) NOT NULL,
      timeStamp DATETIME NOT NULL,
      repoId INT NOT NULL,
      FOREIGN KEY (repoId) REFERENCES vcs.repo(id),
      parentCommitId INT,
      FOREIGN KEY (parentCommitId) REFERENCES vcs.commit(id),
      referencesTree VARCHAR(255),
      FOREIGN KEY (referencesTree) REFERENCES vcs.tree(hash)
    );`,
		`CREATE TABLE IF NOT EXISTS blobContent(
      hash VARCHAR(64) NOT NULL PRIMARY KEY,
      content BLOB NOT NULL
    );`,
		`CREATE TABLE IF NOT EXISTS vcs.tree_entry(
      id INT AUTO_INCREMENT PRIMARY KEY,
      name varchar(255) NOT NULL,
      type INT NOT NULL,
      parentTreeId VARCHAR(64) NOT NULL,
      childBlobId VARCHAR(64),
      childTreeId VARCHAR(64),
      FOREIGN KEY (parentTreeId) REFERENCES vcs.tree(hash),
      FOREIGN KEY (childBlobId) REFERENCES vcs.blobContent(hash),
      FOREIGN KEY (childTreeId) REFERENCES vcs.tree(hash)
    );`,
		`
      DELIMITER //

      CREATE FUNCTION get_latest_commit_id(repoId INT) RETURNS INT
      DETERMINISTIC
      BEGIN
          DECLARE latest_commit INT;
          
          -- Get the latest commit ID for the specified repository
          SELECT id INTO latest_commit
          FROM vcs.commit
          WHERE repoId = repoId
          ORDER BY timeStamp DESC
          LIMIT 1;
          
          -- If no commit found, return NULL
          RETURN latest_commit;
      END //

      DELIMITER;`,
		`
      DELIMITER //

      CREATE TRIGGER before_commit_insert
      BEFORE INSERT ON vcs.commit
      FOR EACH ROW
      BEGIN
          DECLARE latest_commit INT;

          -- Call the function to get the latest commit ID for the repository
          SET latest_commit = get_latest_commit_id(NEW.repoId);
          
          -- Set the parentCommitId in the new row
          SET NEW.parentCommitId = latest_commit;
      END //DELIMITER ;`,
	}

	for _, query := range createTablesQuery {
		_, err = d.Db.Exec(query)

		if err != nil {
			return fmt.Errorf("error creating tables: %s", err)
		}
	}
	return nil
}
