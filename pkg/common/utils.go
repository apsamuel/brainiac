package common

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
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

func GetLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &logger
}
