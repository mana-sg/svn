package user

import (
	"fmt"
	"strings"

	"github.com/mana-sg/vcs/internal/db"
	"github.com/mana-sg/vcs/internal/utils"
	"github.com/mana-sg/vcs/pkg/models"
)

func CreateUser(db db.DbHandler, name string, email string, password string, confirmPass string) error {
	// null value validation
	if strings.Compare(name, "") == 0 {
		return fmt.Errorf("name field cannot be empty")
	}
	if strings.Compare(email, "") == 0 {
		return fmt.Errorf("email field cannot be empty")
	}
	if strings.Compare(password, "") == 0 {
		return fmt.Errorf("password field cannot be empty")
	}
	if strings.Compare(confirmPass, "") == 0 {
		return fmt.Errorf("confirm password field cannot be empty")
	}

	// password and confirm password must be matching or return error
	if strings.Compare(password, confirmPass) != 0 {
		return fmt.Errorf("password and confirm password fields don't match")
	}

	// creating hash of password to maintain security
	hashedPass, err := utils.Hash([]byte(password))
	if err != nil {
		return fmt.Errorf("error in creating hashed password: %v", err)
	}

	// creating and executing query to insert user record into database
	createUserQuery := "INSERT INTO vcs.users (name, email, password) VALUES(?, ?, ?)"
	res, err := db.SetValue(createUserQuery, name, email, hashedPass)
	if err != nil {
		return fmt.Errorf("error inserting the user record into datbase: %v", err)
	}

	// getting user id of the latest insert so that we can create sort of a context for which user is active
	userId, err := res.LastInsertId()

	//  setting the context for current user
	models.SetActiveUser(uint(userId), name)

	return nil
}
