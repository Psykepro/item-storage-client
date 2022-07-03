package _domain

type Publisher interface {
	Publish(body []byte) error
}
