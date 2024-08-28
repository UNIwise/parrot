package project

import (
	"time"

	"github.com/google/uuid"
)

func GenerateUUID() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

func GenerateTimestamp() int64 {
	return time.Now().Unix()
}
