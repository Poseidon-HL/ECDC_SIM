package event_trigger

import (
	"container/heap"
	"testing"
)

func TestEventHeap_Get(t *testing.T) {
	tests := []struct {
		name string
		eq   EventHeap
		want *Event
	}{
		{name: "testEventHeapInitAndGet", eq: *NewEventHeap([]*Event{
			{eventTime: 1},
			{eventTime: 6},
			{eventTime: 2},
		})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.eq.Get()
			t.Log(got.eventTime)
			got = tt.eq.Get()
			t.Log(got.eventTime)
		})
	}
}

func TestNewEventHeap(t *testing.T) {
	type args struct {
		eventList []*Event
	}
	tests := []struct {
		name string
		args args
		want *EventHeap
	}{
		{name: "TestBasicOperations", args: args{eventList: []*Event{
			{eventTime: 9},
		}}},
		{name: "TestBasicOperations", args: args{eventList: []*Event{
			{eventTime: 1},
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEventHeap(tt.args.eventList)
			heap.Push(got, NewEvent(2, 0, 0, 0, nil))
			t.Log(got.Get().eventTime)
		})
	}
}
