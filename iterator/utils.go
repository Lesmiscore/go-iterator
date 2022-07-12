package iterator

// ToSlice converts an Iterator to a slice.
// The given Iterator will have no items remaining and cannot be reused.
func ToSlice[E interface{}](t Iterator[E]) (r []*E, err error) {
	var n *E
	for t.HasNext() {
		n, err = t.GetNext()
		if err != nil {
			return
		}
		r = append(r, n)
	}
	return
}

// Generator will make an Iterator with channel attached.
// This will start a goroutine.
// DO NOT close the given chan inside fn.
// chan will be closed automatically, after fn ends.
func Generator[E interface{}](fn func(chan *E)) Iterator[E] {
	ch := make(chan *E)
	go func() {
		defer close(ch)
		fn(ch)
	}()
	return IteratorFromChan(ch)
}

///////////
/// MAP ///
///////////

type mapIterator[E interface{}, T interface{}] struct {
	src Iterator[E]
	fn  func(int, *E) (*T, error)
	idx int
}

func (t *mapIterator[E, T]) HasNext() bool {
	return t.src.HasNext()
}

func (t *mapIterator[E, T]) GetNext() (e *T, err error) {
	var nxt *E
	nxt, err = t.src.GetNext()
	if err != nil {
		return
	}
	t.idx++
	return t.fn(t.idx, nxt)
}

func Map[E interface{}, T interface{}](src Iterator[E], fn func(int, *E) (*T, error)) Iterator[T] {
	return &mapIterator[E, T]{
		src: src,
		fn:  fn,
		idx: -1,
	}
}

//////////////
/// FILTER ///
//////////////

type filterIterator[E interface{}] struct {
	src       Iterator[E]
	fn        func(int, *E) (bool, error)
	idx       int
	nextState int // [-1, 1] -> [unknown, done, continue]
	nextItem  *E
	lastError error
}

func (t *filterIterator[E]) calcNext() (err error) {
	for t.src.HasNext() {
		var item *E
		item, err = t.GetNext()
		t.idx++
		if err != nil {
			return
		}
		var resp bool
		resp, err = t.fn(t.idx, item)
		if err != nil {
			return
		}
		if resp {
			t.nextItem, t.nextState = item, 1
			return nil
		}
	}
	t.nextState = 0
	return nil
}

func (t *filterIterator[E]) HasNext() bool {
	if t.lastError != nil {
		return false
	}
	if t.nextState == -1 {
		t.lastError = t.calcNext()
		if t.lastError != nil {
			return false
		}
	}
	return t.nextState == 1
}

func (t *filterIterator[E]) GetNext() (e *E, err error) {
	if t.lastError != nil {
		return nil, t.lastError
	}
	if t.nextState == -1 {
		t.lastError = t.calcNext()
		if t.lastError != nil {
			return nil, t.lastError
		}
	}
	e = t.nextItem
	t.nextItem = nil
	t.nextState = -1
	return
}

func Filter[E interface{}](src Iterator[E], fn func(int, *E) (bool, error)) Iterator[E] {
	return &filterIterator[E]{
		src: src,
		fn:  fn,
		idx: -1,
	}
}

// check interfaces
var (
	_ Iterator[int]    = (*mapIterator[string, int])(nil)
	_ Iterator[string] = (*filterIterator[string])(nil)
)
