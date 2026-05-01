package router

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type StringBool struct {
	Str  string
	Flag bool
}

type Channel struct {
	mu                sync.RWMutex
	password          string
	connected_devices map[string]map[string]StringBool
}

func generateToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func setPassword(password string) {
	if password != "" {
		channel.password = hashAndSalt([]byte(password))
	}
}

func verifyPassword(password string) bool {
	if channel.password == "" {
		return true
	}
	return checkPassword(password, channel.password)
}
