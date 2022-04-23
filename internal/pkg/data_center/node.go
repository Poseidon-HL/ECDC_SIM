package data_center

import "ECDC_SIM/internal/pkg/util"

type NodeState int8

const (
	NodeStateNormal NodeState = iota
	NodeStateUnavailable
	NodeStateCrashed
)

type Node struct {
	state                           NodeState
	nodeFailDistribution            *util.Weibull
	nodeTransientFailDistribution   *util.Weibull
	nodeTransientRepairDistribution *util.Weibull
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
			state:                           NodeStateNormal,
			nodeFailDistribution:            nFailD,
			nodeTransientFailDistribution:   nTFailD,
			nodeTransientRepairDistribution: nTRepairD,
		})
	}
	return nodesManager
}
