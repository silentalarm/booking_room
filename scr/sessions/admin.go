package sessions

var adminList = []string{
	"ctristan",
	"hviva",
}

func IsAdmin(nickName string) bool {
	for _, val := range adminList {
		if nickName == val {
			return true
		}
	}
	return false
}
