package pick

import "pick/line"

// Buffer interface is used for containers for lines
type Buffer interface {
	linesInRange(int, int) []line.Line
	LineAt(int) (line.Line, error)
	Size() int
}

// FilteredBuffer holds a "filtered" buffer. It holds a reference to
// the source buffer (note: should be immutable) and a list of indices
// into the source buffer
type FilteredBuffer struct {
	maxcols   int
	src       Buffer
	selection []int // maps from our index to src's index
}

type idgen struct {
	ch chan uint64
}
