package heap

import (
	"container/heap"
)

type Object struct {
	Id       string
	Nodes    []*Object
	Distance int
	Parent   *Object
}

type ObjectHeap []*Object

func NewHeap() *ObjectHeap {
	h := &ObjectHeap{}
	heap.Init(h)
	return h
}

func (h ObjectHeap) Len() int           { return len(h) }
func (h ObjectHeap) Less(i, j int) bool { return h[i].Distance < h[j].Distance }
func (h ObjectHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *ObjectHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Object))
}

func (h *ObjectHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
