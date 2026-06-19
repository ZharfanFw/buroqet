package mocks

type MockStorage struct {
}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (m *MockStorage) Upload(
	fileName string,
) (string, error) {

	return "https://storage.local/" + fileName, nil
}