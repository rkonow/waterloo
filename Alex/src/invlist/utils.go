package invClustering

import (
	"math/rand"
)

type Item struct {
	value    int // The value of the item; arbitrary.
	priority int // The priority of the item in the queue.
	index    int // The index of the item in the heap.
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	a := *pq
	n := len(a)
	a = a[0 : n+1]
	item := x.(*Item)
	item.index = n
	a[n] = item
	*pq = a
}

func (pq *PriorityQueue) Pop() interface{} {
	a := *pq
	n := len(a)
	item := a[n-1]
	item.index = -1 // for safety
	*pq = a[0 : n-1]
	return item
}

//***************** UTILS *******************//
func SwapInt(data []int, i int, j int) []int {
	aux := data[i]
	data[i] = data[j]
	data[j] = aux
	return data
}

func SwapString(data []string, i int, j int) []string {
	aux := data[i]
	data[i] = data[j]
	data[j] = aux
	return data
}

func Shuffle(datas []string, datai []int, dataf []int) ([]string, []int, []int) {
	for i := len(datas) - 1; i > 0; i-- {
		if j := rand.Intn(i + 1); i != j {
			datas = SwapString(datas, i, j)
			datai = SwapInt(datai, i, j)
			dataf = SwapInt(dataf, i, j)
		}
	}
	return datas, datai, dataf
}

//***************** UTILS *******************//
