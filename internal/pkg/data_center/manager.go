package data_center

import "ECDC_SIM/internal/pkg/util"

type DCState int8

const (
	OK DCState = iota
	DEGRADED
)

type DCManager struct {
	state         DCState
	disksManager  *DisksManager
	nodesManager  *NodesManager
	rackManager   *RacksManager
	disksPerNode  int // 每一节点上的磁盘数
	nodesPerRack  int // 每一机架上的节点数
	stripesNum    int
	chunksNum     int
	chunkSize     int
	dataChunksNum int
}

type DCConf struct {
	racksNum                   int
	stripesNum                 int
	disksPerNode               int
	diskCapacity               int
	nodesPerRack               int
	chunkNum                   int
	chunkSize                  int
	dataChunksNum              int
	nFailD, nTFailD, nTRepairD *util.Weibull
	dFailD, dRepairD           *util.Weibull
	rFailD, rRepairD           *util.Weibull
}

func InitDCManager(dcConf *DCConf) *DCManager {
	dcManager := &DCManager{
		state:         OK,
		disksPerNode:  dcConf.disksPerNode,
		nodesPerRack:  dcConf.nodesPerRack,
		stripesNum:    dcConf.stripesNum,
		chunksNum:     dcConf.chunkNum,
		chunkSize:     dcConf.chunkSize,
		dataChunksNum: dcConf.dataChunksNum,
	}
	dcManager.nodesManager = NewNodesManager(dcConf.racksNum*dcConf.nodesPerRack, dcConf.nFailD, dcConf.nTFailD, dcConf.nTRepairD)
	dcManager.disksManager = NewDisksManager(dcManager.nodesManager.nodesNum*dcConf.disksPerNode, dcConf.diskCapacity, dcConf.dFailD, dcConf.dRepairD)
	dcManager.rackManager = NewRacksManager(dcConf.racksNum, dcConf.rFailD, dcConf.rRepairD)
	return dcManager
}

func (dcm *DCManager) GenerateDataPlacement(erConf *ErasureCodeConf) {
	switch erConf.CodeType {
	case RS:

	case LRC:

	default:

	}
}
