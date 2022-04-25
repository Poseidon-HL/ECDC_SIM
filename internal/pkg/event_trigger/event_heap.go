package event_trigger

import (
	"container/heap"
	"github.com/gogap/logrus"
	"reflect"
)

type EventHeap []*Event

func (eq *EventHeap) Len() int {
	return len(*eq)
}

func (eq *EventHeap) Less(i, j int) bool {
	return (*eq)[i].eventTime < (*eq)[j].eventTime
}

func (eq *EventHeap) Swap(i, j int) {
	(*eq)[i], (*eq)[j] = (*eq)[j], (*eq)[i]
}

func (eq *EventHeap) Pop() (v interface{}) {
	*eq, v = (*eq)[:len(*eq)-1], (*eq)[len(*eq)-1]
	return
}

func (eq *EventHeap) Push(v interface{}) {
	*eq = append(*eq, v.(*Event))
}

func NewEventHeap(eventList []*Event) *EventHeap {
	eventQueue := new(EventHeap)
	for _, event := range eventList {
		*eventQueue = append(*eventQueue, event)
	}
	heap.Init(eventQueue)
	return eventQueue
}

func (eq *EventHeap) Get() *Event {
	eventInterface := heap.Pop(eq)
	event, ok := eventInterface.(*Event)
	if !ok {
		logrus.Errorf("[EventHeap.Get] event type error, type=%+v", reflect.TypeOf(event))
	}
	return event
}
