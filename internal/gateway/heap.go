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
	*pq[i].index = j
	*pq[j].index = i
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

type EndpointPriorityQueue struct {
	pq       *priorityQueue
	indexMap map[string]*int
}

func NewEndpointPriorityQueue() *EndpointPriorityQueue {
	pq := make(priorityQueue, 0)
	return &EndpointPriorityQueue{
		pq:       &pq,
		indexMap: make(map[string]*int),
	}
}

func (eq *EndpointPriorityQueue) Size() int {
	return eq.pq.Len()
}

func (eq *EndpointPriorityQueue) Push(address string, priority int) {
	n := eq.pq.Len()
	item := &inactiveEndpointEntry{
		address:  address,
		hitScore: priority,
		index:    &n,
	}
	heap.Push(eq.pq, item)
	eq.indexMap[address] = item.index
}

func (eq *EndpointPriorityQueue) PushNew(address string) {
	n := eq.pq.Len()
	prio := 0
	if n > 0 {
		/*A priority in the middle of the array should be approx. half way down
		  the heap and should hold a median-type value*/
		prio = (*eq.pq)[n/2].hitScore
	}
	eq.Push(address, prio)
}

func (eq *EndpointPriorityQueue) Pop() string {
	if eq.pq.Len() == 0 {
		return ""
	}
	item := heap.Pop(eq.pq)
	endpointEntry := item.(*inactiveEndpointEntry)
	delete(eq.indexMap, endpointEntry.address)

	return endpointEntry.address
}

func (eq *EndpointPriorityQueue) Remove(address string) {
	if index, ok := eq.indexMap[address]; ok {
		heap.Remove(eq.pq, *index)
	}
}
