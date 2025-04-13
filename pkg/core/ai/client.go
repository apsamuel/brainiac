package ai

import (
	"github.com/apsamuel/brainiac/pkg/http"
)

var Client = http.NewApiClient(true, 2000)
