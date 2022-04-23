package event_trigger

type EventType int8

const (
	EventNodeFail EventType = iota
	EventNodeTransientFail
	EventNodeTransientRepair

	EventDiskFail
	EventDiskRepair

	EventRackFail
	EventRackRepair
)

type DeviceType int8

const (
	Rack DeviceType = iota
	Node
	Disk
)

type Event struct {
	eventTime    float64
	eventType    EventType
	deviceIdList []int
	deviceType   DeviceType
	bandwidth    float64
}

type EventManager struct {
	eventQueue *EventHeap
	waitQueue  *EventHeap
}

func (em *EventManager) GetNextEvent() {
	
}
