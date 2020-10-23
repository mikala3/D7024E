package main

type Storage struct {
	str map[string][]byte
}

func NewStorage() *Storage {
	storage := &Storage{}
	m := make(map[string][]byte)
	storage.str = m
	return storage
}

func (storage *Storage) Store(hash string, data string) {
	bytearr := []byte(data)
	storage.str[hash] = bytearr
}

func (storage *Storage) Check(hash string) bool {
	check := storage.str[hash]
	if (len(check) == 0) {
		return false
	}
	return true
}

func (storage *Storage) Get(hash string) []byte {
	check := storage.str[hash]
	return check
}