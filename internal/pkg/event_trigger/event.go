package event_trigger

import (
	"ECDC_SIM/internal/pkg/data_center"
	"container/heap"
	"github.com/gogap/logrus"
)

var EventHandlerFuncMap = map[EventType]EventHandlerFunc{
	EventDiskFail:            DiskFailHandler,
	EventDiskRepair:          DiskRepairHandler,
	EventNodeFail:            NodeFailHandler,
	EventNodeTransientFail:   NodeTransientFailHandler,
	EventNodeTransientRepair: NodeTransientRepairHandler,
	EventRackFail:            RackFailHandler,
	EventRackRepair:          RackRepairHandler,
}

type EventType int8

const (
	EventNodeFail EventType = iota
	EventNodeTransientFail
	EventNodeTransientRepair

	EventDiskFail
	EventDiskRepair

	EventRackFail
	EventRackRepair

	EventMissionEnd
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

type EventExecResult struct {
	EventTime float64
}

func (e *Event) isValidEvent() bool {
	return len(e.deviceIdList) > 0
}

func NewEvent(eventTime float64, eventType EventType, deviceType DeviceType, bandwidth float64, dIdList []int) *Event {
	return &Event{
		eventTime:    eventTime,
		eventType:    eventType,
		deviceIdList: dIdList,
		deviceType:   deviceType,
		bandwidth:    bandwidth,
	}
}

type EventManager struct {
	eventQueue *EventHeap
	waitQueue  *EventHeap

	repairStripesNum            int
	repairStripesSingleChunkNum int
	delayedStripesNum           int
	delayedRepairDict           map[int][]int
}

func NewEventManager() *EventManager {
	return &EventManager{
		eventQueue:        NewEventHeap(make([]*Event, 0)),
		waitQueue:         NewEventHeap(make([]*Event, 0)),
		delayedRepairDict: make(map[int][]int),
	}
}

type EventHandlerFunc func(em *EventManager, event *Event, dList []int, bList []float64) (*Event, error)

func DiskFailHandler(em *EventManager, event *Event, dList []int, bList []float64) (*Event, error) {
	if event.deviceType != Disk {
		logrus.Error("[DiskFailHandler] deviceType wrong")
	}
	failTime := event.eventTime
	diskM := data_center.GetDCManager().DiskManager()
	for _, diskId := range dList {
		if diskM.GetDiskState(diskId) != data_center.DiskStateCrashed {
			if _, ok := em.delayedRepairDict[diskId]; ok {
				delete(em.delayedRepairDict, diskId)
			}
			diskM.FailDisk(diskId, failTime)
			em.SetDiskRepair(diskId, failTime)
		}
	}
	return NewEvent(failTime, EventDiskFail, Disk, 0, dList), nil
}

func DiskRepairHandler(em *EventManager, event *Event, dList []int, bList []float64) (*Event, error) {
	repairTime := event.eventTime
	dcManager := data_center.GetDCManager()
	diskM, nodeM, network := dcManager.DiskManager(), dcManager.NodeManager(), dcManager.Network()
	for _, diskId := range dList {
		if diskM.GetDiskState(diskId) == data_center.DiskStateCrashed {
			diskM.RepairDisk(diskId, repairTime)
			em.SetDiskFail(diskId, repairTime)
		}
		nodeId := dcManager.GetNodeIdByDiskId(diskId)
		if nodeM.GetNodeState(nodeId) == data_center.NodeStateCrashed {
			allDiskOK := true
			for offset := 0; offset < dcManager.GetDisksPerNode(); offset++ {
				if diskM.GetDiskState(dcManager.GetDiskIdByNodeId(nodeId, offset)) != data_center.DiskStateNormal {
					allDiskOK = false
				}
			}
			if allDiskOK {
				nodeM.RepairNode(nodeId)
				if !dcManager.UseTrace() {
					em.SetNodeFail(nodeId, repairTime)
				}
			}
		}
	}
	if network.UseNetwork() {
		for _, bandwidth := range bList {
			network.UpdateAvailCrossRackRepairBandwidth(network.GetAvailCrossRackRepairBandwidth() + bandwidth)
		}
	}
	return NewEvent(repairTime, EventDiskRepair, Disk, 0, dList), nil
}

func NodeFailHandler(em *EventManager, event *Event, dList []int, bList []float64) (*Event, error) {
	failTime := event.eventTime
	failedDiskList := make([]int, 0)
	dcManager := data_center.GetDCManager()
	diskM, nodeM := dcManager.DiskManager(), dcManager.NodeManager()
	for _, nodeId := range dList {
		if nodeM.GetNodeState(nodeId) != data_center.NodeStateCrashed {
			nodeM.FailNode(nodeId, failTime)
		}
		for offset := 0; offset < dcManager.GetDisksPerNode(); offset++ {
			diskId := dcManager.GetDiskIdByNodeId(nodeId, offset)
			failedDiskList = append(failedDiskList, diskId)
			if diskM.GetDiskState(diskId) != data_center.DiskStateCrashed {
				// TODO here some questions
				if _, ok := em.delayedRepairDict[diskId]; ok {
					delete(em.delayedRepairDict, diskId)
				}
				diskM.FailDisk(diskId, failTime)
				em.SetDiskRepair(diskId, failTime)
			}
		}
	}
	return NewEvent(failTime, EventNodeFail, Disk, 0, failedDiskList), nil
}

func NodeTransientFailHandler(em *EventManager, event *Event, dList []int, bList []float64) (*Event, error) {
	failTime := event.eventTime
	dcManager := data_center.GetDCManager()
	diskM, nodeM := dcManager.DiskManager(), dcManager.NodeManager()
	for _, nodeId := range dList {
		if nodeM.GetNodeState(nodeId) == data_center.NodeStateNormal {
			nodeM.OfflineNode(nodeId)
			for offset := 0; offset < dcManager.GetDisksPerNode(); offset++ {
				diskId := dcManager.GetDiskIdByNodeId(nodeId, offset)
				if diskM.GetDiskState(diskId) == data_center.DiskStateNormal {
					diskM.OfflineDisk(diskId, failTime)
				}
			}
		}
		if dcManager.UseTrace() {
			em.SetNodeTransientRepair(nodeId, failTime)
		}
	}
	return NewEvent(failTime, EventNodeTransientFail, Node, 0, nil), nil
}

func NodeTransientRepairHandler(em *EventManager, event *Event, dList []int, bList []float64) (*Event, error) {
	dcManager := data_center.GetDCManager()
	diskM, nodeM := dcManager.DiskManager(), dcManager.NodeManager()
	repairTime := event.eventTime
	for _, nodeId := range dList {
		if nodeM.GetNodeState(nodeId) == data_center.NodeStateUnavailable {
			nodeM.OnlineNode(nodeId)
			for offset := 0; offset < dcManager.GetDisksPerNode(); offset++ {
				diskId := dcManager.GetDiskIdByNodeId(nodeId, offset)
				if diskM.GetDiskState(diskId) == data_center.DiskStateUnavailable {
					diskM.OnlineDisk(diskId, repairTime)
				}
			}
		}
		if !dcManager.UseTrace() {
			em.SetNodeTransientFail(nodeId, repairTime)
		}
	}
	return NewEvent(repairTime, EventNodeTransientRepair, Node, 0, nil), nil
}

func RackFailHandler(em *EventManager, event *Event, dList []int, bList []float64) (*Event, error) {
	failTime := event.eventTime
	dcManager := data_center.GetDCManager()
	diskM, nodeM, rackM := dcManager.DiskManager(), dcManager.NodeManager(), dcManager.RackManager()
	for _, rackId := range dList {
		if rackM.GetRackState(rackId) == data_center.RackStateNormal {
			rackM.FailRack(rackId)
			for offset := 0; offset < dcManager.GetNodesPerRack(); offset++ {
				nodeId := dcManager.GetNodeIdByRackId(rackId, offset)
				if nodeM.GetNodeState(nodeId) == data_center.NodeStateNormal {
					nodeM.OfflineNode(nodeId)
					for diskOffset := 0; diskOffset < dcManager.GetDisksPerNode(); diskOffset++ {
						diskId := dcManager.GetDiskIdByNodeId(nodeId, diskOffset)
						if diskM.GetDiskState(diskId) == data_center.DiskStateNormal {
							diskM.OfflineDisk(diskId, failTime)
						}
					}
				}
			}
		}
		if !dcManager.UsePowerOutage() {
			em.SetRackRepair(rackId, failTime)
		}
	}
	return NewEvent(failTime, EventRackFail, Rack, 0, nil), nil
}

func RackRepairHandler(em *EventManager, event *Event, dList []int, bList []float64) (*Event, error) {
	repairTime := event.eventTime
	dcManager := data_center.GetDCManager()
	diskM, nodeM, rackM := dcManager.DiskManager(), dcManager.NodeManager(), dcManager.RackManager()
	for _, rackId := range dList {
		if rackM.GetRackState(rackId) == data_center.RackStateUnavailable {
			rackM.RepairRack(rackId)
			for offset := 0; offset < dcManager.GetNodesPerRack(); offset++ {
				nodeId := dcManager.GetNodeIdByRackId(rackId, offset)
				if nodeM.GetNodeState(nodeId) == data_center.NodeStateUnavailable {
					nodeM.OnlineNode(nodeId)
					for diskOffset := 0; diskOffset < dcManager.GetDisksPerNode(); diskOffset++ {
						diskId := dcManager.GetDiskIdByNodeId(nodeId, diskOffset)
						if diskM.GetDiskState(diskId) == data_center.DiskStateUnavailable {
							diskM.OnlineDisk(diskId, repairTime)
						}
					}
				}
			}
		}
		if !dcManager.UsePowerOutage() {
			em.SetRackFail(rackId, repairTime)
		}
	}
	return NewEvent(repairTime, EventRackFail, Rack, 0, nil), nil
}

// HandleNextEvent 根据事件队列进行相应的事件操作
func (em *EventManager) HandleNextEvent(currentTime float64) *EventExecResult {
	var err error
	dcManager := data_center.GetDCManager()
	em.checkDelayedRepairDict()
	em.checkWaitQueue(currentTime)
	event := em.eventQueue.Get()
	deviceList, repairBandwidthList := em.popSameEvent(event)
	if event.eventTime > dcManager.GetMissionTime() {
		return &EventExecResult{EventTime: currentTime}
	}
	if handleFunc, ok := EventHandlerFuncMap[event.eventType]; ok {
		event, err = handleFunc(em, event, deviceList, repairBandwidthList)
		if err != nil {
			logrus.Error("[EventManager.GetNextEvent] EventHandlerFuncMap error")
		}
		return &EventExecResult{EventTime: event.eventTime}
	} else {
		logrus.Error("[EventManager.GetNextEvent] HandlerFunc missing")
	}

	return nil
}

func (em *EventManager) checkDelayedRepairDict() {
	if len(em.delayedRepairDict) == 0 {
		return
	}
	diskToRemove := make([]int, 0)
	dcManager := data_center.GetDCManager()
	diskManager := dcManager.DiskManager()
	for diskIdKey, stripeIdList := range em.delayedRepairDict {
		newDictValue := make([]int, 0)
		for _, stripeId := range stripeIdList {
			var repairDelay bool
			var numUnavailingChunk int
			for _, diskId := range dcManager.GetStripesLocation(stripeId) {
				if diskManager.GetDiskState(diskId) != data_center.DiskStateNormal {
					numUnavailingChunk++
				}
				if numUnavailingChunk > dcManager.ErasureCodeConf().N-dcManager.ErasureCodeConf().K {
					repairDelay = true
					break
				}
			}
			if repairDelay {
				newDictValue = append(newDictValue, stripeId)
			}
		}
		if len(newDictValue) == 0 {
			diskToRemove = append(diskToRemove, diskIdKey)
		} else {
			em.delayedRepairDict[diskIdKey] = newDictValue
		}
	}
	for _, key := range diskToRemove {
		delete(em.delayedRepairDict, key)
	}
}

func (em *EventManager) checkWaitQueue(currentTime float64) {
	if len(*em.waitQueue) == 0 {
		return
	}
	dcManager := data_center.GetDCManager()
	networkM := dcManager.Network()
	rackManager := dcManager.RackManager()
	diskId := (*em.waitQueue)[0].deviceIdList[0]
	rackId := dcManager.GetRackIdByDiskId(diskId)
	if networkM.UseNetwork() && networkM.GetAvailCrossRackRepairBandwidth() != 0 &&
		networkM.GetAvailIntraRackRepairBandwidth(rackId) != 0 &&
		rackManager.GetRackState(rackId) == data_center.RackStateNormal {
		heap.Pop(em.waitQueue)
		em.SetDiskRepair(diskId, currentTime)
	}
}

func (em *EventManager) popSameEvent(event *Event) ([]int, []float64) {
	dcManager := data_center.GetDCManager()
	networkM := dcManager.Network()
	deviceIdList := make([]int, 0)
	deviceIdList = append(deviceIdList, event.deviceIdList...)
	repairBandwidthList := make([]float64, 0)
	if networkM.UseNetwork() && event.eventType == EventDiskRepair {
		repairBandwidthList = append(repairBandwidthList, event.bandwidth)
	}
	for (*em.eventQueue)[0].eventTime == event.eventTime &&
		(*em.eventQueue)[0].eventType == event.eventType {
		event = em.eventQueue.Get()
		deviceIdList = append(deviceIdList, event.deviceIdList...)
		if networkM.UseNetwork() && event.eventType == EventDiskRepair {
			repairBandwidthList = append(repairBandwidthList, event.bandwidth)
		}
	}
	return deviceIdList, repairBandwidthList
}

func (em *EventManager) SetDiskRepair(diskId int, currentTime float64) {
	dcManager := data_center.GetDCManager()
	networkM := dcManager.Network()
	rackM := dcManager.RackManager()
	diskM := dcManager.DiskManager()
	rackId := dcManager.GetRackIdByDiskId(diskId)
	if networkM.GetAvailCrossRackRepairBandwidth() == 0 || rackM.GetRackState(rackId) != data_center.RackStateNormal {
		heap.Push(em.waitQueue, NewEvent(currentTime, EventDiskFail, Disk, 0, []int{diskId}))
		return
	}
	crossRackDownload := 0
	// TODO check logic here
	stripeIdList := diskM.GetDiskStripes(diskId)
	em.repairStripesNum += len(stripeIdList)
	var stripesToDelay []int
	// 针对这一个块上的所有条带，均需要进行修复
	for _, stripeId := range stripeIdList {
		numOfFailedChunks, numOfAliveChunkInSameRack, numOfUnavailingChunk := 0, 0, 0
		for _, diskNum := range dcManager.GetStripesLocation(stripeId) {
			if diskM.GetDiskState(diskNum) != data_center.DiskStateNormal {
				numOfUnavailingChunk++
			}
			switch dcManager.ErasureCodeConf().CodeType {
			case data_center.RS:
				if diskM.GetDiskState(diskNum) == data_center.DiskStateCrashed {
					numOfFailedChunks++
				} else if dcManager.GetRackIdByDiskId(diskId) == rackId {
					numOfAliveChunkInSameRack++
				}
			case data_center.LRC:
				// TODO
			}
		}
		if numOfFailedChunks == 1 {
			em.repairStripesSingleChunkNum++
		}
		// 无法完成纠删码要求的修复
		if numOfUnavailingChunk > (dcManager.ErasureCodeConf().N - dcManager.ErasureCodeConf().K) {
			stripesToDelay = append(stripesToDelay, stripeId)
		}
		switch dcManager.ErasureCodeConf().CodeType {
		case data_center.RS:
			if numOfAliveChunkInSameRack < dcManager.ErasureCodeConf().K {
				crossRackDownload += dcManager.ErasureCodeConf().K - numOfAliveChunkInSameRack
			}
		case data_center.LRC:
		}
	}
	repairBandwidth := networkM.GetAvailCrossRackRepairBandwidth()
	networkM.UpdateAvailCrossRackRepairBandwidth(0)
	repairTime := float64(crossRackDownload*dcManager.GetChunkSize()) / repairBandwidth
	repairTime /= float64(3600)
	if len(stripesToDelay) > 0 {
		em.delayedStripesNum += len(stripesToDelay)
		em.delayedRepairDict[diskId] = stripesToDelay
	}

	// TODO repair bandwidth
	heap.Push(em.eventQueue, NewEvent(repairTime+currentTime, EventDiskRepair, Disk, repairBandwidth, []int{diskId}))
}

func (em *EventManager) SetDiskFail(diskId int, currentTime float64) {
	dcManager := data_center.GetDCManager()
	diskM := dcManager.DiskManager()
	heap.Push(em.eventQueue, NewEvent(diskM.GetDiskFailDistribution(diskId).Draw()+currentTime,
		EventDiskFail, Disk, 0, []int{diskId}))
}

func (em *EventManager) SetNodeTransientRepair(nodeId int, currentTime float64) {
	dcManager := data_center.GetDCManager()
	nodeM := dcManager.NodeManager()
	heap.Push(em.eventQueue, NewEvent(nodeM.GetTransitRepairDistribution(nodeId).Draw()+currentTime,
		EventNodeTransientRepair, Node, 0, []int{nodeId}))
}

func (em *EventManager) SetNodeFail(nodeId int, currentTime float64) {
	dcManager := data_center.GetDCManager()
	nodeM := dcManager.NodeManager()
	heap.Push(em.eventQueue, NewEvent(nodeM.GetNodeFailDistribution(nodeId).Draw()+currentTime,
		EventNodeFail, Node, 0, []int{nodeId}))
}

func (em *EventManager) SetNodeTransientFail(nodeId int, currentTime float64) {
	dcManager := data_center.GetDCManager()
	nodeM := dcManager.NodeManager()
	heap.Push(em.eventQueue, NewEvent(nodeM.GetTransitFailDistribution(nodeId).Draw()+currentTime,
		EventNodeFail, Node, 0, []int{nodeId}))
}

func (em *EventManager) SetRackRepair(rackId int, currentTime float64) {
	dcManager := data_center.GetDCManager()
	rackM := dcManager.RackManager()
	heap.Push(em.eventQueue, NewEvent(rackM.GetRackRepairDistribution(rackId).Draw()+currentTime,
		EventRackRepair, Rack, 0, []int{rackId}))
}

func (em *EventManager) SetRackFail(rackId int, currentTime float64) {
	dcManager := data_center.GetDCManager()
	rackM := dcManager.RackManager()
	heap.Push(em.eventQueue, NewEvent(rackM.GetRackFailDistribution(rackId).Draw()+currentTime,
		EventRackRepair, Rack, 0, []int{rackId}))
}
