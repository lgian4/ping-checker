package duration_circular_buffer

import "time"

type Duration_CircularBuffer struct {
	data            []time.Duration
	head, tail      int
	isFull, isEmpty bool
	Size            int
	Length          int
}

func New(size int) *Duration_CircularBuffer {
	return &Duration_CircularBuffer{
		data:    make([]time.Duration, size),
		isFull:  false,
		isEmpty: true,
		Size:    size,
	}
}

func (cb *Duration_CircularBuffer) Enqueue(value time.Duration) error {
	if cb.isFull {
		cb.head = (cb.head + 1) % len(cb.data)
	} else {
		cb.Length += 1
	}

	cb.data[cb.tail] = value
	cb.tail = (cb.tail + 1) % len(cb.data)
	cb.isEmpty = false

	if !cb.isFull && cb.tail == cb.head {
		cb.isFull = true
	}

	return nil
}

func (cb *Duration_CircularBuffer) Get(index int) time.Duration {
	if cb.isFull {
		index = (cb.head + index) % len(cb.data)
		return cb.data[index]
	} else if index < cb.Length {
		return cb.data[index]
	}
	return -1
}
