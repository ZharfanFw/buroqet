package mocks

type MockKafka struct {
}

func NewMockKafka() *MockKafka {
	return &MockKafka{}
}

func (m *MockKafka) Publish(
	topic string,
	message string,
) error {

	return nil
}