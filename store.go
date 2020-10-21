package main

type Storage struct {
	str map[string][]byte
}

func NewStorage(str map[string][]byte) *Storage {
	storage := &Storage{}
	storage.str = str
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