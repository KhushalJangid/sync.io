package router

type StringBool struct {
	Str  string
	Flag bool
}

type Channel struct {
	password          string
	connected_devices map[string]map[string]StringBool
}

func setPassword(password string) {
	if password != "" {
		channel.password = hashAndSalt([]byte(password))
	}
}
func verifyPassword(password string) bool {
	hash := hashAndSalt([]byte(password))
	if hash == channel.password {
		return false
	} else {
		return true
	}
}
