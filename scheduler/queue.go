package main

import (
	"context"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/samber/lo"
	"sync"
	"time"
)

type queue struct {
	QueuedIDs     []string `json:"queuedIds"`
	ProcessingIDs []string `json:"processingIds"`
	concurrency   int
	processID     func(id string) error
}

func newQueue(concurrency int, processID func(id string) error) *queue {
	return &queue{
		concurrency: concurrency,
		processID:   processID,
	}
}

func (q *queue) full() bool {
	return len(q.ProcessingIDs) >= q.concurrency
}

func (q *queue) doCycle() {
	for len(q.QueuedIDs) > 0 && !q.full() {
		var id string
		id, q.QueuedIDs = q.QueuedIDs[0], q.QueuedIDs[1:]

		q.add(id)
	}
}

func (q *queue) next() {
	if len(q.QueuedIDs) > 0 && !q.full() {
		var id string
		id, q.QueuedIDs = q.QueuedIDs[0], q.QueuedIDs[1:]

		q.add(id)
	}

	if len(q.QueuedIDs) == 0 && len(q.ProcessingIDs) == 0 {
		log.L.Debug().Msg("queue finished")
	}
}

func (q *queue) queue(id string) {
	q.QueuedIDs = append(q.QueuedIDs, id)
}

var sliceLock = sync.Mutex{}

func (q *queue) add(id string) {
	sliceLock.Lock()
	defer sliceLock.Unlock()

	q.ProcessingIDs = append(q.ProcessingIDs, id)

	go func() {
		log.L.Debug().Str("id", id).Msg("processing id")
		err := q.processID(id)
		if err != nil {
			log.L.Error().Err(err).Send()
			return
		}
		q.remove(id)
		q.next()
		log.L.Debug().Str("id", id).Msg("processed id")
	}()
}

func (q *queue) remove(id string) {
	sliceLock.Lock()
	defer sliceLock.Unlock()

	q.ProcessingIDs = lo.Filter(q.ProcessingIDs, func(i string, _ int) bool {
		return i != id
	})
}

func (q *queue) run(ctx context.Context) {
	q.doCycle()
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			q.doCycle()
		case <-ctx.Done():
			return
		}
	}
}
