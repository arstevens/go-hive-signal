package analyzer

import (
	"container/list"
	"sync"
)

var FrequencyAveragingWidth = 50

type swarmTracker struct {
	historyLength          int
	frequencyHistory       *list.List
	fHistoryMutex          *sync.Mutex
	dspaceFrequencyHistory map[string]*list.List
	dfHistoryMutex         *sync.Mutex
}

func newTracker() *swarmTracker {
	return &swarmTracker{
		historyLength:          FrequencyAveragingWidth,
		dspaceFrequencyHistory: make(map[string]*list.List),
		dfHistoryMutex:         &sync.Mutex{},
	}
}

func (st *swarmTracker) AddFrequencyDatapoint(dspace string, record int) {
	st.dfHistoryMutex.Lock()
	datapoints := st.dspaceFrequencyHistory[dspace]
	if datapoints == nil {
		datapoints = initHistoryQueue(st.historyLength)
		st.dspaceFrequencyHistory[dspace] = datapoints
	}

	datapoints.PushBack(record)
	datapoints.Remove(datapoints.Front())
	st.dfHistoryMutex.Unlock()
}

func initHistoryQueue(size int) *list.List {
	queue := list.New()
	for i := 0; i < size; i++ {
		queue.PushBack(1)
	}
	return queue
}

func (st *swarmTracker) Cleanup() {
	dspaceHistory := st.CalculateDataspaceFrequencies()
	for dspace, freq := range dspaceHistory {
		if freq <= 0 {
			st.dfHistoryMutex.Lock()
			delete(st.dspaceFrequencyHistory, dspace)
			st.dfHistoryMutex.Unlock()
		}
	}
}

func (st *swarmTracker) CalculateFrequency() int {
	dspaceFrequencies := st.CalculateDataspaceFrequencies()
	frequency := 0
	for _, freq := range dspaceFrequencies {
		frequency += freq
	}
	return frequency
}

func (st *swarmTracker) CalculateDataspaceFrequencies() map[string]int {
	st.dfHistoryMutex.Lock()
	dspaceFrequencies := make(map[string]int)
	for dspace, freqList := range st.dspaceFrequencyHistory {
		dspaceFrequencies[dspace] = averageFrequencies(freqList)
	}
	st.dfHistoryMutex.Unlock()
	return dspaceFrequencies
}

func averageFrequencies(datapoints *list.List) int {
	listLen := datapoints.Len()
	if listLen == 0 {
		return 0
	}

	freq := 0
	for e := datapoints.Front(); e != nil; e = e.Next() {
		freq += e.Value.(int)
	}
	return freq / listLen
}
