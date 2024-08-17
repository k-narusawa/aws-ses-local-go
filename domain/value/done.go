package value

type Done struct {
	ValueObject[bool]
}

func NewDone(done bool) Done {
	return Done{ValueObject: ValueObject[bool]{value: done}}
}

func (d Done) Bool() bool {
	return d.value
}

func OfDone(value bool) Done {
	return Done{ValueObject: ValueObject[bool]{value: value}}
}
