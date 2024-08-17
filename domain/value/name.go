package value

type Name struct {
	ValueObject[string]
}

func NewName(name string) Name {
	return Name{
		ValueObject: ValueObject[string]{value: name},
	}
}
