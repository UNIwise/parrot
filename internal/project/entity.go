package project

import "time"

type Version struct {
	ID         string
	Name       string
	CreatedAt  time.Time
}

type Project struct {
	ID               int64
	Name             string
	NumberOfVersions int
	CreatedAt        string
}
