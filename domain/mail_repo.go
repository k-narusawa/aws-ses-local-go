package domain

type IMailRepo interface {
	Store(Mail) error
	Delete(string) error
	DeleteAll() error
}
