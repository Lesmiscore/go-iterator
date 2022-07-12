package iterator

type chanIterator[E interface{}] struct {
	ch      chan *E
	hasNext bool
	buffer  *E
}

func IteratorFromChan[E interface{}](ch chan *E) Iterator[E] {
	return &chanIterator[E]{
		ch:      ch,
		hasNext: false,
		buffer:  nil,
	}
}

func (t *chanIterator[E]) HasNext() bool {
	if t.hasNext {
		return true
	}
	t.buffer, t.hasNext = <-t.ch
	return t.hasNext
}

func (t *chanIterator[E]) GetNext() (*E, error) {
	if t.hasNext {
		t.hasNext = false
		defer func() {
			t.buffer = nil
		}()
		return t.buffer, nil
	}
	return <-t.ch, nil
}

var (
	_ Iterator[int] = (*chanIterator[int])(nil)
)
