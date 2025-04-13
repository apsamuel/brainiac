package common

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func getAppRoot() string {
	// This is a placeholder function. You should implement the logic to get the module path.
	cwd, _ := os.Getwd()
	dir := filepath.Clean(cwd)
	for {
		if file, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !file.IsDir() {
			return dir
		}
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		dir = d
	}
	return ""
}

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
