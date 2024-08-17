package domain

type IMailRepo interface {
	Store(Mail) error
}
