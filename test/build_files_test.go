package test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/mana-sg/vcs/internal/db"
	"github.com/mana-sg/vcs/internal/repository"
)

func TestBuildFileTree(t *testing.T) {
	godotenv.Load("../.env")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	dbname := os.Getenv("DATABASE_NAME")

	dbcf := db.DbHandler{}

	err := dbcf.ConfigDB(user, password, host, dbname)
	if err != nil {
		t.Errorf("Error connecting to database: %v", err)
	}

	// err = dbcf.PrepareDb()
	// if err != nil {
	// 	t.Errorf("error preparing the database: %v", err)
	// }

	files, err := repository.GetAllFilesForCommit(dbcf, 2)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Log(files)
}
