package gateway

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestEndpointPriorityQueue(t *testing.T) {
	fmt.Printf("---------------------\nINACTIVE PRIORITY QUEUE TEST\n---------------------------\n")
	pq := NewEndpointPriorityQueue()

	rand.Seed(time.Now().UnixNano())
	totalItems := 50
	for i := 0; i < totalItems; i++ {
		prio := rand.Intn(100)
		addr := strconv.Itoa(prio)
		pq.Push(addr, prio)
	}

	fmt.Printf("Random Insert Pop Order: \n")
	printPQ(pq)

	for i := 0; i < totalItems; i++ {
		addr := strconv.Itoa(i)
		pq.PushNew(addr)
	}
	fmt.Printf("\nPushNew Insertions Pop Order: \n")
	printPQ(pq)

	removal := strconv.Itoa(totalItems / 2)
	pq.PushNew(removal)
	fmt.Printf("\nRemoving Entry %s: \n", removal)
	pq.Remove(removal)
	printPQ(pq)

}

func printPQ(pq *EndpointPriorityQueue) {
	n := pq.Size()
	if n == 0 {
		fmt.Printf("\tEMPTY\n")
		return
	}
	for i := 0; i < n; i++ {
		if i%10 == 0 {
			fmt.Printf("\n\t")
		}
		fmt.Printf("%s ", pq.Pop())
	}
}
