package tracker

import (
	"container/list"
	"sync"
)

var FrequencyAveragingWidth = 50

type SwarmLoadTracker struct {
	frequencyHistory *list.List
	historyMutex     *sync.Mutex
}

func NewLoadTracker(historyCap int) *SwarmLoadTracker {
	if historyCap <= 0 {
		historyCap = FrequencyAveragingWidth
	}

	return &SwarmLoadTracker{
		frequencyHistory: initHistoryQueue(historyCap),
		historyMutex:     &sync.Mutex{},
	}
}

func (st *SwarmLoadTracker) AddFrequencyDatapoint(record int) {
	st.historyMutex.Lock()
	st.frequencyHistory.PushBack(record)
	st.frequencyHistory.Remove(st.frequencyHistory.Front())
	st.historyMutex.Unlock()
}

func (st *SwarmLoadTracker) CalculateAverageFrequency() int {
	freq := 0
	for e := st.frequencyHistory.Front(); e != nil; e = e.Next() {
		freq += e.Value.(int)
	}
	return freq / st.frequencyHistory.Len()
}

func initHistoryQueue(size int) *list.List {
	queue := list.New()
	for i := 0; i < size; i++ {
		queue.PushBack(1)
	}
	return queue
}
