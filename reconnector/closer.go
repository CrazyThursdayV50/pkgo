package reconnector

type Closer interface {
	Close()
}

type emptyCloser struct{}

func (emptyCloser) Close() {}

var _ Closer = emptyCloser{}

type ErrorCloser interface {
	Close() error
}

type WrappedErrorCloser[C Closer] struct {
	Conn C
}

func (w *WrappedErrorCloser[C]) Close() error {
	w.Conn.Close()
	return nil
}

var _ ErrorCloser = (*WrappedErrorCloser[emptyCloser])(nil)

func NewWrappedErrorCloser[C Closer](c C) *WrappedErrorCloser[C] {
	var w WrappedErrorCloser[C]
	w.Conn = c
	return &w
}

type ClosedChecker interface {
	Closed() bool
}

type CloserClosedChecker interface {
	Closer
	ClosedChecker
}

type ErrorCloserClosedChecker interface {
	ErrorCloser
	ClosedChecker
}

type emptyClosedChecker struct{}

func (emptyClosedChecker) Closed() bool {
	return false
}

type WrappedErrorCloserChecker[C ErrorCloser] struct {
	emptyClosedChecker
	Conn C
}

func (w *WrappedErrorCloserChecker[C]) Close() error {
	return w.Conn.Close()
}

var _ ErrorCloserClosedChecker = (*WrappedErrorCloserChecker[*WrappedErrorCloser[emptyCloser]])(nil)

func NewWrappedErrorCloserChecker[C ErrorCloser](c C) *WrappedErrorCloserChecker[C] {
	var w WrappedErrorCloserChecker[C]
	w.Conn = c
	return &w
}
