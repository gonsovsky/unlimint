package test

import (
	"../shared"
	"../storage"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

//Симулятор отправки сообщений ...
type Siumulator struct {
	config    shared.TestCfg
	buffer    storage.IRepository
	db        storage.IRepository
	sentToWeb int32
}

//Start send fake messages to web server
func NewTest(config shared.TestCfg, buffer storage.IRepository, db storage.IRepository) *Siumulator {
	e := Siumulator{config: config, buffer: buffer, db: db}
	time.Sleep(300 * time.Millisecond)
	go e.Stats()
	for i := 1; i <= e.config.NumberOfTestHits; i++ {
		var msg = shared.GoogleHit{
			ClientID:        fmt.Sprintf("%03d", i),
			TrackingID:      "UA-146615186-1",
			DocumentPath:    fmt.Sprintf("%03d", i),
			HitType:         "pageview",
			ProtocolVersion: "1",
		}
		e.post(&msg)
		time.Sleep(10 * time.Millisecond)
	}
	return &e
}

// Send fake message to web host
func (e *Siumulator) post(hit *shared.GoogleHit) error {
	defer atomic.AddInt32(&e.sentToWeb, 1)
	g := NewApi(e.config.ApiUrl)
	return g.Send(hit)
}

func (c *Siumulator) GetCount() int32 {
	return atomic.LoadInt32(&c.sentToWeb)
}

func (e *Siumulator) Stats() {
	for range time.NewTicker(1500 * time.Millisecond).C {
		fmt.Printf("отправлено %d, в буфере: %d, в базе: %d\r\n",
			e.GetCount(), e.buffer.GetCount(), e.db.GetCount())
		if e.db.GetCount() >= int32(e.config.NumberOfTestHits) {
			os.Exit(0)
		}

	}
}
