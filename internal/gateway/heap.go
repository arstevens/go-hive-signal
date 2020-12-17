package gateway

import (
	"container/heap"
)

type inactiveEndpointEntry struct {
	address  string
	hitScore int
	index    *int
}

type priorityQueue []*inactiveEndpointEntry

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].hitScore < pq[j].hitScore
}
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	*pq[j].index = j
	*pq[i].index = i
}
func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*inactiveEndpointEntry)
	*item.index = n
	*pq = append(*pq, item)
}
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*item.index = -1
	*pq = old[0 : n-1]
	return item
}

type endpointPriorityQueue struct {
	pq       *priorityQueue
	indexMap map[string]*int
}

func newEndpointPriorityQueue() *endpointPriorityQueue {
	pq := make(priorityQueue, 0)
	return &endpointPriorityQueue{
		pq:       &pq,
		indexMap: make(map[string]*int),
	}
}

func (eq *endpointPriorityQueue) IsEmpty() bool {
	return eq.GetSize() == 0
}

func (eq *endpointPriorityQueue) GetSize() int {
	return eq.pq.Len()
}

func (eq *endpointPriorityQueue) Push(address string, priority int) {
	n := eq.pq.Len()
	item := &inactiveEndpointEntry{
		address:  address,
		hitScore: priority,
		index:    &n,
	}
	heap.Push(eq.pq, item)
	eq.indexMap[address] = item.index
}

func (eq *endpointPriorityQueue) PushNew(address string) {
	n := eq.pq.Len()
	prio := 0
	if n > 0 {
		/*A priority in the middle of the array should be approx. half way down
		  the heap and should hold a median-type value*/
		prio = (*eq.pq)[n/2].hitScore
	}
	eq.Push(address, prio)
}

func (eq *endpointPriorityQueue) Pop() (string, int) {
	if eq.pq.Len() == 0 {
		return "", 0
	}
	item := heap.Pop(eq.pq)
	endpointEntry := item.(*inactiveEndpointEntry)
	delete(eq.indexMap, endpointEntry.address)

	return endpointEntry.address, endpointEntry.hitScore
}

func (eq *endpointPriorityQueue) Remove(address string) {
	if index, ok := eq.indexMap[address]; ok {
		heap.Remove(eq.pq, *index)
	}
}
