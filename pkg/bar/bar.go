package bar

import (
	"time"

	"github.com/cheggaaa/pb"
)

const RefreshRate = time.Millisecond * 100

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type WriteCounter struct {
	n   int // bytes read so far
	bar *pb.ProgressBar
}

func NewWriteCounter(total int, filepath string) *WriteCounter {
	b := pb.New(total)
	b.SetRefreshRate(RefreshRate)
	b.ShowTimeLeft = false
	b.ShowSpeed = true
	b.ShowCounters = true
	b.SetUnits(pb.U_BYTES)
	b.Prefix(filepath)

	return &WriteCounter{
		bar: b,
	}
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	wc.n += len(p)
	wc.bar.Set(wc.n)
	return wc.n, nil
}

func (wc *WriteCounter) Start() {
	wc.bar.Start()
}

func (wc *WriteCounter) Finish() {
	wc.bar.Finish()
}
