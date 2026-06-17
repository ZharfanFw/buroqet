package storage

type LocalStorage struct {
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (s *LocalStorage) Upload(
	fileName string,
) (string, error) {

	return "https://storage.local/" + fileName, nil
}