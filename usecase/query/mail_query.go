package query

type IMailDtoQueryService interface {
	FindAll(limit, offset int) ([]MailDto, error)
	FindByMessageID(messageID string) (MailDto, error)
	FindByTo(to *string, limit, offset int) ([]MailDto, error)
}
