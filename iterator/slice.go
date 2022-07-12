package iterator

type sliceIterator[E interface{}] struct {
	slc []*E
}

func IteratorFromSlice[E interface{}](s []*E) Iterator[E] {
	return &sliceIterator[E]{
		slc: s,
	}
}

func (t *sliceIterator[E]) HasNext() bool {
	for range t.slc {
		return true
	}
	return false
}

func (t *sliceIterator[E]) GetNext() (e *E, _ error) {
	for _, e = range t.slc {
		t.slc = t.slc[1:]
		break
	}
	return
}

var (
	_ Iterator[int] = (*sliceIterator[int])(nil)
)
