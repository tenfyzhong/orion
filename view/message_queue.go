package view

import (
	"github.com/tenfyzhong/orion/model"
)

// MessageQueue queue to maintain message
type MessageQueue struct {
	head     int
	tail     int
	messages []*model.Message
}

// NewMessageQueue Create message queue
func NewMessageQueue(capture int) *MessageQueue {
	return &MessageQueue{
		head:     0,
		tail:     0,
		messages: make([]*model.Message, capture+1, capture+1),
	}
}

// Empty the queue is empty
func (mq MessageQueue) Empty() bool {
	return mq.head == mq.tail
}

// Capacity the queue's capacity
func (mq MessageQueue) Capacity() int {
	return len(mq.messages) - 1
}

// Full the queue is full
func (mq MessageQueue) Full() bool {
	return (mq.tail+1)%len(mq.messages) == mq.head
}

// Push push message into queue, return true if success
func (mq *MessageQueue) Push(m *model.Message) bool {
	if mq.Full() {
		return false
	}
	mq.messages[mq.tail] = m
	mq.tail = (mq.tail + 1) % len(mq.messages)

	return true
}

// Pop pop message, return nil if the queue is empty
func (mq *MessageQueue) Pop() *model.Message {
	if mq.Empty() {
		return nil
	}
	m := mq.messages[mq.head]
	mq.head = (mq.head + 1) % len(mq.messages)
	return m
}

// Head get the head node
func (mq MessageQueue) Head() *model.Message {
	if mq.Empty() {
		return nil
	}
	return mq.messages[mq.head]
}

// SearchByNum get the message by message.num
func (mq MessageQueue) SearchByNum(num uint32) *model.Message {
	i := mq.head
	for {
		if i == mq.tail {
			return nil
		}
		if mq.messages[i].Num == num {
			return mq.messages[i]
		}
		i = (i + 1) % len(mq.messages)
	}
}
