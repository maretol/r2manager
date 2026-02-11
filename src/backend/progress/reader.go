package progress

import (
	"errors"
	"io"
	"sync/atomic"
	"time"

	serviceif "r2manager/service/interface"
)

var ErrNilCallback = errors.New("callback cannot be nil")

// ProgressReadCloser は io.ReadCloser をラップして読み取りバイト数を追跡する。
// Phase 1（クライアントからのリクエストボディ受信）で使用する。
type ProgressReadCloser struct {
	inner     io.ReadCloser
	processed atomic.Int64
	callback  serviceif.ProgressCallback
	lastEmit  time.Time
	throttle  time.Duration
}

func NewProgressReadCloser(inner io.ReadCloser, callback serviceif.ProgressCallback, throttle time.Duration) (*ProgressReadCloser, error) {
	if callback == nil {
		return nil, ErrNilCallback
	}
	return &ProgressReadCloser{
		inner:    inner,
		callback: callback,
		throttle: throttle,
	}, nil
}

func (r *ProgressReadCloser) Read(p []byte) (int, error) {
	n, err := r.inner.Read(p)
	if n > 0 {
		processed := r.processed.Add(int64(n))
		now := time.Now()
		if now.Sub(r.lastEmit) >= r.throttle || err == io.EOF {
			r.lastEmit = now
			r.callback(processed)
		}
	}
	return n, err
}

func (r *ProgressReadCloser) Close() error {
	r.callback(r.processed.Load())
	return r.inner.Close()
}

// ProgressReadSeeker は io.ReadSeeker をラップして読み取りバイト数を追跡する。
// Phase 2（R2への PutObject アップロード）で使用する。
type ProgressReadSeeker struct {
	inner     io.ReadSeeker
	processed atomic.Int64
	callback  serviceif.ProgressCallback
	lastEmit  time.Time
	throttle  time.Duration
}

func NewProgressReadSeeker(inner io.ReadSeeker, callback serviceif.ProgressCallback, throttle time.Duration) (*ProgressReadSeeker, error) {
	if callback == nil {
		return nil, ErrNilCallback
	}
	return &ProgressReadSeeker{
		inner:    inner,
		callback: callback,
		throttle: throttle,
	}, nil
}

func (r *ProgressReadSeeker) Read(p []byte) (int, error) {
	n, err := r.inner.Read(p)
	if n > 0 {
		processed := r.processed.Add(int64(n))
		now := time.Now()
		if now.Sub(r.lastEmit) >= r.throttle || err == io.EOF {
			r.lastEmit = now
			r.callback(processed)
		}
	}
	return n, err
}

func (r *ProgressReadSeeker) Seek(offset int64, whence int) (int64, error) {
	pos, err := r.inner.Seek(offset, whence)
	if err == nil && offset == 0 && whence == io.SeekStart {
		r.processed.Store(0)
	}
	return pos, err
}
