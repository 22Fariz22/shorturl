package producer_consumer

import "os"

type producer struct {
	file *os.File // файл для записи
}

func NewProducer(filename string) (*producer, error) {
	// открываем файл для записи в конец
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &producer{file: file}, nil
}

func (p *producer) Close() error {
	// закрываем файл
	return p.file.Close()
}

type consumer struct {
	file *os.File // файл для чтения
}

func NewConsumer(filename string) (*consumer, error) {
	// открываем файл для чтения
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &consumer{file: file}, nil
}

func (c *consumer) Close() error {
	// закрываем файл
	return c.file.Close()
}
