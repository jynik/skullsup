// SPDX License Identifier: MIT
package server

import (
	"errors"
	"sync"

	"github.com/jynik/skullsup/go/src/network"
)

type queue []network.Message

type MessageQueue struct {
	mutex  sync.Mutex
	queues map[string]queue

	maxQueues int
	maxDepth  int
}

func NewMessageQueue(maxQueues, maxDepth int) *MessageQueue {
	var q MessageQueue
	q.queues = make(map[string]queue)
	q.maxQueues = maxQueues
	q.maxDepth = maxDepth
	return &q
}

func (mq *MessageQueue) Enqueue(key string, m network.Message) error {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	q, exists := mq.queues[key]
	if (!exists && len(mq.queues) >= mq.maxQueues) || len(q) > mq.maxDepth {
		return errors.New(network.ErrorQueueFull)
	}

	if exists {
		mq.queues[key] = append(q, m)
	} else {
		mq.queues[key] = []network.Message{m}
	}

	return nil
}

func (mq *MessageQueue) Dequeue(key string) (network.Message, error) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	q, exists := mq.queues[key]
	if !exists || len(q) == 0 {
		return network.Message{}, errors.New(network.ErrorQueueEmpty)
	}

	msg := q[0]
	mq.queues[key] = q[1:]
	if len(mq.queues[key]) == 0 {
		delete(mq.queues, key)
	}

	return msg, nil
}
