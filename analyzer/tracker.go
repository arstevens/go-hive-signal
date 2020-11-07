package analyzer

import "container/list"

type swarmTracker struct {
	historyLength          int
	frequencyHistory       *list.List
	dspaceFrequencyHistory map[string]*list.List
}

func (st *swarmTracker) AddFrequencyDatapoint(record int) {
	st.frequencyHistory.PushBack(record)
	st.frequencyHistory.Remove(st.frequencyHistory.Front())
}

func (st *swarmTracker) AddDataspaceFrequencyDatapoint(dspace string, record int) {
	datapoints := st.dspaceFrequencyHistory[dspace]
	if datapoints == nil {
		datapoints = initHistoryQueue(st.historyLength)
		st.dspaceFrequencyHistory[dspace] = datapoints
	}

	datapoints.PushBack(record)
	datapoints.Remove(datapoints.Front())
}

func initHistoryQueue(size int) *list.List {
	queue := list.New()
	for i := 0; i < size; i++ {
		queue.PushBack(0)
	}
	return queue
}

func (st *swarmTracker) CalculateFrequency() int {
	return averageFrequencies(st.frequencyHistory)
}

func (st *swarmTracker) CalculateDataspaceFrequencies() map[string]int {
	dspaceFrequencies := make(map[string]int)
	for dspace, freqList := range st.dspaceFrequencyHistory {
		dspaceFrequencies[dspace] = averageFrequencies(freqList)
	}
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
