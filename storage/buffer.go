package storage

import (
	"../shared"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//Временный буффер (на 5 секунд) в памяти программы
type Buffer struct {
	PersistentDb  IRepository
	flushInterval int
	flushItems    int
	count         int32
	hits          []shared.GoogleHit
	m             sync.RWMutex
}

func NewBuffer(persistDb IRepository, cfg shared.SetupCfg) *Buffer {
	b := Buffer{PersistentDb: persistDb, flushInterval: cfg.FlushInterval, flushItems: cfg.FlushItems}
	return &b
}

func (b *Buffer) Consume() {
	for range time.NewTicker(time.Duration(b.flushInterval) * time.Second).C {
		b.Flush()
	}
}

func (b *Buffer) Post(hit shared.GoogleHit) error {
	b.m.Lock()
	defer b.m.Unlock()
	b.hits = append(b.hits, hit)
	b.count++
	return nil
}

func (b *Buffer) Flush() error {
	b.m.Lock()
	x := (len(b.hits))
	if x > b.flushItems {
		x = b.flushItems
	}
	tmp := b.hits[:x]
	b.hits = b.hits[x:]
	b.m.Unlock()
	if x == 0 {
		return nil
	}
	fmt.Println("сохраняем ", x, " хитов...")
	for _, hit := range tmp {
		b.PersistentDb.Post(hit)
		atomic.AddInt32(&b.count, int32(-1))
		time.Sleep(1 * time.Millisecond)
	}
	return nil
}

func (b *Buffer) GetCount() int32 {
	return atomic.LoadInt32(&b.count)
}
