package query

type ICountQueryService interface {
	CountByTo(to *string) (int, error)
}
