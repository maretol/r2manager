package domain

import "time"

type Object struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
}

type ListObjectsResult struct {
	Objects                 []Object `json:"objects"`
	Prefix                  string   `json:"prefix"`
	Delimiter               string   `json:"delimiter"`
	IsTruncated             bool     `json:"is_truncated"`
	NextContinuationToken   string   `json:"next_continuation_token,omitempty"`
}
