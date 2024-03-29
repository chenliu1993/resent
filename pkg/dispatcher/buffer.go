package dispatcher

import (
	"sync"
)

type buffer struct {
	m     sync.RWMutex
	data  []interface{}
	avail int
	cap   int
}

func NewBuffer(n int) *buffer {
	return &buffer{
		data:  make([]interface{}, n, n),
		avail: 0,
		cap:   n,
	}
}

func (b *buffer) ReadOne() interface{} {
	b.m.RLock()
	defer b.m.RUnlock()

	if b.avail == 0 {
		return nil
	}

	data := b.data[0]
	b.data = b.data[1:]
	b.data = append(b.data, nil)
	b.avail--

	return data
}

func (b *buffer) WriteOne(data interface{}) {
	b.m.Lock()
	defer b.m.Unlock()
	if b.avail < b.cap {
		b.data[b.avail] = data
		b.avail++
		return
	}

	newData := make([]interface{}, 2*b.avail)
	copy(newData, b.data)
	b.data = newData
	b.data[b.avail] = data
	b.cap = 2 * b.avail
	b.avail = b.avail + 1
}
