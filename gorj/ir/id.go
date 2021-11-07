package ir

type ID int

type idAlloc struct {
	nextID ID
}

func (ia *idAlloc) next() ID {
	v := ia.nextID
	ia.nextID++
	return v
}

func (ia *idAlloc) count() ID {
	return ia.nextID + 1
}
