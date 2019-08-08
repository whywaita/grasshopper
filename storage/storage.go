package storage

type Storage interface {
	Put(string) error
}
