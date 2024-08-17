package value

type Title struct {
	ValueObject[string]
}

func NewTitle(title string) Title {
	return Title{ValueObject: ValueObject[string]{value: title}}
}
