package user

import (
	"fmt"
	"strings"

	"github.com/mana-sg/vcs/types"
	"github.com/mana-sg/vcs/utils"
)

func CreateUser(db utils.DbHandler, name string, email string, password string, confirmPass string) error {
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
	types.ChooseUser(uint(userId), name)

	return nil
}

func LogIn(db utils.DbHandler, email string, password string) error {
	if strings.Compare(email, "") == 0 {
		return fmt.Errorf("email field cannot be empty")
	}
	if strings.Compare(password, "") == 0 {
		return fmt.Errorf("password field cannot be empty")
	}

	hashedPass, err := utils.Hash([]byte(password))
	if err != nil {
		return fmt.Errorf("error in creating hashed password: %v", err)
	}

	createUserQuery := "SELECT id, name, password FROM vcs.users where email= ?"

	rows, err := db.GetValue(createUserQuery, email)
	if err != nil {
		return fmt.Errorf("error inserting the user record into datbase: %v", err)
	}
	defer rows.Close()

	var userId uint
	var name, passwordRes string

	if err = rows.Scan(&userId, &name, &passwordRes); err != nil {
		return fmt.Errorf("error retrieving data: %v", err)
	}

	if strings.Compare(string(hashedPass), passwordRes) != 0 {
		return fmt.Errorf("Wrong Password! Please try again!")
	}

	types.ChooseUser(uint(userId), name)

	return nil
}
