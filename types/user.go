// Used to instantiate which user is live when the program is running
// Can check which user using this curr user variable
// Can get the name of the curr user using this file

package types

type user struct {
	id   uint
	name string
}

var currUser user

func ChooseUser(id uint, name string) {
	currUser.id = id
	currUser.name = name
}

func GetCurrUser() string {
	return currUser.name
}
