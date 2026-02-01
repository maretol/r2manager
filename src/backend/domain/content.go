package domain

import "io"

type ObjectContent struct {
	Body        io.ReadCloser
	ContentType string
	Size        int64
	ETag        string
	CacheHit    bool
}
