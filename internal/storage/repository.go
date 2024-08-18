//go:generate mockgen --source=repository.go -destination=repository_mock.go -package=storage
package storage

type Storage interface {
	// TODO: Add methods to implement
}
