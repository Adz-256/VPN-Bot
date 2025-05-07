package broker

type Consumer interface {
	Read(key string) ([]byte, error)
}

type Publisher interface {
	Write(key string, value []byte) error
}
