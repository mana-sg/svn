// Used to instantiate which user is live when the program is running
// Can check which user using this curr user variable
// Can get the name of the curr user using this file

package models

type user struct {
	id   uint
	name string
}

var activeUser user

func SetActiveUser(id uint, name string) {
	activeUser.id = id
	activeUser.name = name
}

func GetActiveUser() (uint, string) {
	return activeUser.id, activeUser.name
}
