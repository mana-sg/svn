package db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type DbHandler struct {
	db *sql.DB
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
	d.db = db

	return nil
}

func (d *DbHandler) Close() error {
	err := d.db.Close()
	if err != nil {
		return fmt.Errorf("error closing database: %v", err)
	}
	return nil
}

func (d *DbHandler) SetValue(query string, args ...interface{}) (sql.Result, error) {
	res, err := d.db.Exec(query, args)
	if err != nil {
		return nil, fmt.Errorf("error inserting values: %v", err)
	}
	return res, nil
}

func (d *DbHandler) GetValue(query string, args ...interface{}) (*sql.Rows, error) {
	res, err := d.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting the value: %v", err)
	}
	return res, nil
}

func (d *DbHandler) PrepareDb() error {
	createDbQuery := "CREATE DATABASE IF NOT EXISTS vcs"

	_, err := d.db.Exec(createDbQuery)

	if err != nil {
		return fmt.Errorf("error creating database: %s", err)
	}

	_, err = d.db.Exec("USE vcs")
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
		`CREATE TABLE IF NOT EXISTS commit(
      id INT AUTO_INCREMENT PRIMARY KEY,
      message VARCHAR(255) NOT NULL,
      timeStamp DATETIME NOT NULL
      repoId INT NOT NULL,
      FOREIGN KEY (repoId) REFERENCES vcs.repo(id)
    );`,
		`CREATE TABLE IF NOT EXISTS repo(
      id INT AUTO_INCREMENT PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      timeCreation DATETIME NOT NULL 
      userId INT,
      FOREIGN KEY (userId) REFERENCES vcs.users(id)
    );`,
		`CREATE TABLE IF NOT EXISTS tree(
      hash VARCHAR(64) NOT NULL PRIMARY KEY 
      tree_entry VARCHAR(64) NOT NULL,
      FOREIGN KEY (tree_entry) references vcs.tree_entry(hash)
    );`,
		`CREATE TABLE IF NOT EXISTS blobContent(
      hash VARCHAR(64) NOT NULL PRIMARY KEY,
      content BLOB NOT NULL
    );`,
		`CREATE TABLE IF NOT EXISTS tree_entry(
      id INT AUTO_INCREMENT PRIMARY KEY,
      name varchar(255) NOT NULL,
      type INT NOT NULL,
      
    );`,
	}

	for _, query := range createTablesQuery {
		_, err = d.db.Exec(query)

		if err != nil {
			return fmt.Errorf("error creating tables: %s", err)
		}
	}
	return nil
}
