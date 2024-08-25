package query

import "aws-ses-local-go/internal/gateways/dao"

type CountQueryService struct {
	MailDao dao.MailDao
}

func NewCountQueryService(
	mailDao dao.MailDao,
) *CountQueryService {
	return &CountQueryService{
		MailDao: mailDao,
	}
}

func (s *CountQueryService) CountByTo(to *string) (int, error) {
	count, err := s.MailDao.CountByTo(to)
	if err != nil {
		return 0, err
	}

	return count, nil
}
