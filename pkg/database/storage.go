package database

import (
	"github.com/apsamuel/brainiac/pkg/common"
)

type Storage struct {
	Name          string
	Type          string
	TrainingStore common.Storer[TrainingDataSchema]
}
