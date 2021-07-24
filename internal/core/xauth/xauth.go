package xauth

import "../xconfig"

func Authenticate(username, password string) bool {
	un := "User_" + username
	passHash, _ := xconfig.ConfigGet(un, "")
	if passHash == password {
		return true
	} else {
		return false
	}
}
