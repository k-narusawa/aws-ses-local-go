package mail

import "aws-ses-local-go/domain"

type Service struct {
	MailRepo domain.IMailRepo
}

func NewService(
	mailRepo domain.IMailRepo,
) *Service {
	return &Service{
		MailRepo: mailRepo,
	}
}

func (s *Service) DeleteMail(mId string) error {
	err := s.MailRepo.Delete(mId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteMails() error {
	err := s.MailRepo.DeleteAll()
	if err != nil {
		return err
	}
	return nil
}
