// SPDX License Identifier: MIT
package server

import (
	"errors"
	"sync"

	c "../common"
)

type queue []c.Message

type MessageQueue struct {
	mutex  sync.Mutex
	queues map[string]queue

	maxKeys     int
	maxElements int
}

func NewMessageQueue(maxKeys, maxLength int) *MessageQueue {
	var q MessageQueue
	q.queues = make(map[string]queue)
	q.maxKeys = maxKeys
	q.maxElements = maxLength
	return &q
}

func (mq *MessageQueue) Enqueue(key string, m c.Message) error {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	q, exists := mq.queues[key]
	if (!exists && len(mq.queues) >= mq.maxKeys) || len(q) > mq.maxElements {
		return errors.New(c.ERR_FULL)
	}

	if exists {
		mq.queues[key] = append(q, m)
	} else {
		mq.queues[key] = []c.Message{m}
	}

	return nil
}

func (mq *MessageQueue) Dequeue(key string) (c.Message, error) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	q, exists := mq.queues[key]
	if !exists || len(q) == 0 {
		return c.Message{}, errors.New(c.ERR_EMPTY)
	}

	msg := q[0]
	mq.queues[key] = q[1:]
	if len(mq.queues[key]) == 0 {
		delete(mq.queues, key)
	}

	return msg, nil
}
