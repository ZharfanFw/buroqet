package kafka

type Producer struct {
}

func NewProducer() *Producer {
	return &Producer{}
}

func (p *Producer) Publish(
	topic string,
	message string,
) error {

	return nil
}