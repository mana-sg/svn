package user

import (
	"fmt"
	"strings"

	"github.com/mana-sg/vcs/internal/db"
	"github.com/mana-sg/vcs/internal/utils"
	"github.com/mana-sg/vcs/pkg/models"
)

func LogIn(db db.DbHandler, email string, password string) error {
	// field validation
	if strings.Compare(email, "") == 0 {
		return fmt.Errorf("email field cannot be empty")
	}
	if strings.Compare(password, "") == 0 {
		return fmt.Errorf("password field cannot be empty")
	}

	// calculating hash password to compare
	hashedPass, err := utils.Hash([]byte(password))
	if err != nil {
		return fmt.Errorf("error in creating hashed password: %v", err)
	}

	//creating the select query to check if the user exists
	createUserQuery := "SELECT id, name, password FROM vcs.users where email= ?"

	rows, err := db.GetValue(createUserQuery, email)
	if err != nil {
		return fmt.Errorf("error getting the user record into datbase: %v", err)
	}
	defer rows.Close()

	if rows == nil {
		return fmt.Errorf("user does not exist")
	}

	// getting the value of the user if the user exists
	var userId uint
	var name, passwordRes string

	rows.Next()
	if err = rows.Scan(&userId, &name, &passwordRes); err != nil {
		return fmt.Errorf("error retrieving data: %v", err)
	}

	// checking if the user entered password and original password match
	if strings.Compare(string(hashedPass), passwordRes) != 0 {
		return fmt.Errorf("Wrong Password! Please try again!")
	}

	// setting the user context for the logged in user
	models.SetActiveUser(uint(userId), name)

	return nil
}
