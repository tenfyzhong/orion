package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tenfyzhong/orion/model"
)

func TestEmpty(t *testing.T) {
	mq := NewMessageQueue(10)
	assert.True(t, mq.Empty())
}

func TestCapacity(t *testing.T) {
	mq := NewMessageQueue(10)
	assert.Equal(t, 10, mq.Capacity())
	assert.Equal(t, 11, len(mq.messages))
}

func TestPush(t *testing.T) {
	mq := NewMessageQueue(2)
	m := &model.Message{}
	assert.True(t, mq.Push(m))
	assert.False(t, mq.Empty())
	assert.Equal(t, 0, mq.head)
	assert.Equal(t, 1, mq.tail)
	assert.Equal(t, m, mq.messages[0])
	assert.True(t, mq.Push(m))
	assert.False(t, mq.Push(m))
	assert.True(t, mq.Full())
}

func TestPop(t *testing.T) {
	mq := NewMessageQueue(2)
	assert.Nil(t, mq.Pop())
	m := &model.Message{}
	mq.Push(m)
	assert.Equal(t, m, mq.Pop())
	mq.Empty()
}

func TestFull(t *testing.T) {
	mq := NewMessageQueue(2)
	m := &model.Message{}
	mq.Push(m)
	assert.False(t, mq.Full())
	mq.Push(m)
	assert.True(t, mq.Full())
}

func TestHead(t *testing.T) {
	mq := NewMessageQueue(2)
	assert.Nil(t, mq.Head())
	m := &model.Message{}
	mq.Push(m)
	assert.Equal(t, m, mq.Head())
}

func TestSearchByNum(t *testing.T) {
	mq := NewMessageQueue(5)
	assert.Nil(t, mq.SearchByNum(1))
	m1 := &model.Message{
		Num: 1,
	}
	mq.Push(m1)
	assert.Equal(t, m1, mq.SearchByNum(m1.Num))
	assert.Nil(t, mq.SearchByNum(2))
	mq.Pop()
	assert.Nil(t, mq.SearchByNum(1))
}

func TestSearchByNum2(t *testing.T) {
	mq := NewMessageQueue(5)

	m1 := &model.Message{
		Num: 1,
	}
	m2 := &model.Message{
		Num: 2,
	}
	m3 := &model.Message{
		Num: 3,
	}
	m4 := &model.Message{
		Num: 4,
	}
	m5 := &model.Message{
		Num: 5,
	}
	m6 := &model.Message{
		Num: 6,
	}
	mq.Push(m1)
	mq.Push(m2)
	mq.Push(m3)
	mq.Push(m4)
	mq.Push(m5)
	assert.True(t, mq.Full())
	mq.Pop()
	mq.Pop()
	mq.Push(m6)
	assert.Equal(t, 0, mq.tail)
	assert.Equal(t, m4, mq.SearchByNum(m4.Num))
}
