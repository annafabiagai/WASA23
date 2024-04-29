package api

import (
	"regexp"
	"strings"

	"github.com/annafabia03/WASA23/service/database"
)

type User struct {
	ID       uint64 `json:"IDuser"`
	Nickname string `json:"nickname"`
}

func (u *User) HasValidUsername() bool {
	nickname := strings.TrimSpace(u.Nickname)
	if nickname == "" {
		return false
	}
	if len(nickname) < 3 || len(nickname) > 16 {
		return false
	}
	match, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_]{2,15}$", nickname)
	return match
}

func (u *User) ToDatabase() database.User {
	return database.User{
		ID:       u.ID,
		Nickname: u.Nickname,
	}
}

func (u *User) FromDatabase(user database.User) {
	u.ID = user.ID
	u.Nickname = user.Nickname
}
