package utils

import (
	"crypto/rand"
	"encoding/hex"
	"meatsrv/constants"
)

func GenerateSecureToken() string {
	return GenerateRandomString(constants.TokenLength)
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
