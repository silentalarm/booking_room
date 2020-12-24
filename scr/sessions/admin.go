package sessions

import (
	"os"
	"strings"
)

var adminList = strings.Split(os.Getenv("ADMINS"), ",")

func IsAdmin(nickName string) bool {
	for _, val := range adminList {
		if nickName == val {
			return true
		}
	}
	return false
}
