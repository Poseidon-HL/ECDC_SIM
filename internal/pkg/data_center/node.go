package data_center

import "ECDC_SIM/internal/pkg/util"

type NodeState int8

const (
	NodeStateNormal NodeState = iota
	NodeStateUnavailable
	NodeStateCrashed
	NodeStateUndefined
)

type Node struct {
	nodeClock                       *DeviceClock
	state                           NodeState
	nodeFailDistribution            *util.Weibull
	nodeTransientFailDistribution   *util.Weibull
	nodeTransientRepairDistribution *util.Weibull
}

func (n *Node) ResetState() {
	n.state = NodeStateNormal
}

func (n *Node) GetState() NodeState {
	return n.state
}

func (n *Node) Fail(currentTime float64) {
	n.state = NodeStateUnavailable
	n.nodeClock.repairTime = 0
	n.nodeClock.repairStart = currentTime
}

func (n *Node) Repair() {
	n.state = NodeStateNormal
	n.nodeClock.globalTime = n.nodeClock.lastUpdateTime
	n.nodeClock.localTime = 0
	n.nodeClock.repairTime = 0
}

func (n *Node) Offline() {
	if n.state == NodeStateNormal {
		n.state = NodeStateUnavailable
	}
}

func (n *Node) Online() {
	if n.state == NodeStateUnavailable {
		n.state = NodeStateNormal
	}
}

type NodesManager struct {
	nodesNum       int
	nodes          []*Node
	failedNodesNum int
	failedNodesMap map[int]int
}

func NewNodesManager(nodesNum int, nFailD, nTFailD, nTRepairD *util.Weibull) *NodesManager {
	nodesManager := &NodesManager{
		nodesNum:       nodesNum,
		failedNodesNum: 0,
		failedNodesMap: make(map[int]int),
	}
	for i := 0; i < nodesNum; i++ {
		nodesManager.nodes = append(nodesManager.nodes, &Node{
			nodeClock:                       new(DeviceClock),
			state:                           NodeStateNormal,
			nodeFailDistribution:            nFailD,
			nodeTransientFailDistribution:   nTFailD,
			nodeTransientRepairDistribution: nTRepairD,
		})
	}
	return nodesManager
}

func (nm *NodesManager) Reset(currentTime float64) {
	for _, node := range nm.nodes {
		node.ResetState()
	}
	nm.failedNodesNum = 0
	nm.failedNodesMap = make(map[int]int)
}

func (nm *NodesManager) isValidNodeId(nodeId int) bool {
	return nodeId >= 0 && nodeId < len(nm.nodes)
}

func (nm *NodesManager) GetNodeState(nodeId int) NodeState {
	if nm.isValidNodeId(nodeId) {
		return nm.nodes[nodeId].GetState()
	}
	return NodeStateUndefined
}

func (nm *NodesManager) FailNode(nodeId int, currentTime float64) {
	if nm.isValidNodeId(nodeId) {
		nm.nodes[nodeId].Fail(currentTime)
	}
}

func (nm *NodesManager) RepairNode(nodeId int) {
	if nm.isValidNodeId(nodeId) {
		nm.nodes[nodeId].Repair()
	}
}

func (nm *NodesManager) OnlineNode(nodeId int) {
	if nm.isValidNodeId(nodeId) {
		nm.nodes[nodeId].Online()
	}
}

func (nm *NodesManager) OfflineNode(nodeId int) {
	if nm.isValidNodeId(nodeId) {
		nm.nodes[nodeId].Offline()
	}
}

func (nm *NodesManager) GetTransitRepairDistribution(nodeId int) *util.Weibull {
	if nm.isValidNodeId(nodeId) {
		return nm.nodes[nodeId].nodeTransientRepairDistribution
	}
	return nil
}

func (nm *NodesManager) GetTransitFailDistribution(nodeId int) *util.Weibull {
	if nm.isValidNodeId(nodeId) {
		return nm.nodes[nodeId].nodeTransientFailDistribution
	}
	return nil
}

func (nm *NodesManager) GetNodeFailDistribution(nodeId int) *util.Weibull {
	if nm.isValidNodeId(nodeId) {
		return nm.nodes[nodeId].nodeFailDistribution
	}
	return nil
}

func (nm *NodesManager) GetNodeNum() int {
	return len(nm.nodes)
}
