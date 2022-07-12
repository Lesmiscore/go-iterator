package iterator

type Iterator[E interface{}] interface {
    HasNext() bool
    GetNext() (*E, error)
}
