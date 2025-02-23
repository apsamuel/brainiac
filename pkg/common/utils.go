package common

import (
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
)

func GetRandomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return fmt.Sprintf("%X", b)[:n]
}

func GetUUID() string {
	id := uuid.New()
	return id.String()
}
