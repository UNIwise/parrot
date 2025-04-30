package project

import (
	"time"

	"github.com/google/uuid"
	"github.com/uniwise/parrot/pkg/poedit"
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

func GetContentMetaMap() map[string]poedit.ContentMeta {
	return poedit.ContentMetaMap
}
