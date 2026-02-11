package progress

import (
	"context"
	"log"
	"sync"
	"time"

	"r2manager/domain"
)

const (
	entryTTL        = 5 * time.Minute
	cleanupInterval = 1 * time.Minute
	channelBuffer   = 32
)

type uploadEntry struct {
	uploadID    string
	createdAt   time.Time
	completedAt *time.Time
	lastEvent   *domain.UploadEvent
	subscribers []chan domain.UploadEvent
	mu          sync.Mutex
}

type UploadProgressStore struct {
	mu      sync.RWMutex
	entries map[string]*uploadEntry
}

func NewUploadProgressStore() *UploadProgressStore {
	return &UploadProgressStore{
		entries: make(map[string]*uploadEntry),
	}
}

// StartCleanupLoop は完了済みエントリを定期的に削除するゴルーチンを起動する。
func (s *UploadProgressStore) StartCleanupLoop(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.cleanup()
			}
		}
	}()
}

func (s *UploadProgressStore) cleanup() {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, entry := range s.entries {
		entry.mu.Lock()
		shouldDelete := false
		if entry.completedAt != nil && now.Sub(*entry.completedAt) > entryTTL {
			shouldDelete = true
			for _, ch := range entry.subscribers {
				close(ch)
			}
			entry.subscribers = nil
		}
		entry.mu.Unlock()

		if shouldDelete {
			delete(s.entries, id)
			log.Printf("cleaned up upload progress entry: %s", id)
		}
	}
}

// Register は新しいアップロードのエントリを作成する。
func (s *UploadProgressStore) Register(uploadID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[uploadID] = &uploadEntry{
		uploadID:  uploadID,
		createdAt: time.Now(),
	}
}

// Publish は進捗イベントをストアに記録し、全subscriberに配信する。
func (s *UploadProgressStore) Publish(uploadID string, event domain.UploadEvent) {
	s.mu.RLock()
	entry, ok := s.entries[uploadID]
	s.mu.RUnlock()
	if !ok {
		return
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	entry.lastEvent = &event
	if event.EventType == domain.EventComplete || event.EventType == domain.EventError {
		now := time.Now()
		entry.completedAt = &now
	}

	for _, ch := range entry.subscribers {
		select {
		case ch <- event:
		default:
		}
	}
}

// Subscribe は指定uploadIDの進捗イベントチャネルを返す。
// 既に完了済みの場合、lastEventを含むチャネルを返してすぐ閉じる。
func (s *UploadProgressStore) Subscribe(uploadID string) (<-chan domain.UploadEvent, func()) {
	ch := make(chan domain.UploadEvent, channelBuffer)

	s.mu.RLock()
	entry, ok := s.entries[uploadID]
	s.mu.RUnlock()

	if !ok {
		close(ch)
		return ch, func() {}
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if entry.completedAt != nil && entry.lastEvent != nil {
		ch <- *entry.lastEvent
		close(ch)
		return ch, func() {}
	}

	entry.subscribers = append(entry.subscribers, ch)

	unsubscribe := func() {
		entry.mu.Lock()
		defer entry.mu.Unlock()
		for i, sub := range entry.subscribers {
			if sub == ch {
				entry.subscribers = append(entry.subscribers[:i], entry.subscribers[i+1:]...)
				break
			}
		}
	}

	return ch, unsubscribe
}
