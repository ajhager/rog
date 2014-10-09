package rog

import "container/heap"

type Actor interface {
	Act() float64
}

type Queue struct {
	aq actorQueue
}

func NewQueue() *Queue {
	aq := make(actorQueue, 0)
	heap.Init(&aq)
	return &Queue{aq}
}

func (q *Queue) Push(actor Actor, time float64) {
	heap.Push(&q.aq, &timedActor{actor: actor, time: time})
}

func (q *Queue) Pop() Actor {
	if q.aq.Len() == 0 {
		return nil
	}
	chosen := heap.Pop(&q.aq).(*timedActor)
	for _, ta := range q.aq {
		ta.time -= chosen.time
	}
	return chosen.actor
}

func (q *Queue) Len() int {
	return q.aq.Len()
}

type timedActor struct {
	index int
	time  float64
	actor Actor
}

type actorQueue []*timedActor

func (aq actorQueue) Len() int {
	return len(aq)
}

func (aq actorQueue) Less(i, j int) bool {
	return aq[i].time < aq[j].time
}

func (aq actorQueue) Swap(i, j int) {
	aq[i], aq[j] = aq[j], aq[i]
	aq[i].index = i
	aq[j].index = j
}

func (aq *actorQueue) Push(x interface{}) {
	n := len(*aq)
	ta := x.(*timedActor)
	ta.index = n
	*aq = append(*aq, ta)
}

func (aq *actorQueue) Pop() interface{} {
	old := *aq
	n := len(old)
	ta := old[n-1]
	ta.index = -1
	*aq = old[0 : n-1]
	return ta
}
