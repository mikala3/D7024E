package main

type Storage struct {
	str map[string]string
}

func NewStorage() *Storage {
	storage := &Storage{}
	m := make(map[string]string)
	storage.str = m
	return storage
}

func (storage *Storage) Store(hash string, data string) {
	storage.str[hash] = data
}

func (storage *Storage) Check(hash string) bool {
	check := storage.str[hash]
	if (len(check) == 0) {
		return false
	}
	return true
}

func (storage *Storage) Get(hash string) string {
	check := storage.str[hash]
	return check
}