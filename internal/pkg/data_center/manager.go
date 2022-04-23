package data_center

import (
	"ECDC_SIM/internal/pkg/enum_error"
	"ECDC_SIM/internal/pkg/util"
	"github.com/gogap/logrus"
)

type DCState int8

const (
	OK DCState = iota
	DEGRADED
)

type DCManager struct {
	state           DCState
	disksManager    *DisksManager
	nodesManager    *NodesManager
	rackManager     *RacksManager
	disksPerNode    int // 每一节点上的磁盘数
	nodesPerRack    int // 每一机架上的节点数
	stripesNum      int
	chunksNum       int
	chunkSize       int
	dataChunksNum   int
	erasureCodeConf *ErasureCodeConf
	stripesLocation [][]int
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

func InitDCManager(dcConf *DCConf, eCConf *ErasureCodeConf) *DCManager {
	dcManager := &DCManager{
		state:           OK,
		disksPerNode:    dcConf.disksPerNode,
		nodesPerRack:    dcConf.nodesPerRack,
		stripesNum:      dcConf.stripesNum,
		chunksNum:       dcConf.chunkNum,
		chunkSize:       dcConf.chunkSize,
		dataChunksNum:   dcConf.dataChunksNum,
		erasureCodeConf: eCConf,
	}
	dcManager.nodesManager = NewNodesManager(dcConf.racksNum*dcConf.nodesPerRack, dcConf.nFailD, dcConf.nTFailD, dcConf.nTRepairD)
	dcManager.disksManager = NewDisksManager(dcManager.nodesManager.nodesNum*dcConf.disksPerNode, dcConf.diskCapacity, dcConf.dFailD, dcConf.dRepairD)
	dcManager.rackManager = NewRacksManager(dcConf.racksNum, dcConf.rFailD, dcConf.rRepairD)
	return dcManager
}

func (dcm *DCManager) GenerateDataPlacement() {
	var err error
	switch dcm.erasureCodeConf.CodeType {
	case RS:
		err = dcm.GeneratePlacementByArchType()
		if err != nil {
			logrus.Errorf("DCManager.GenerateDataPlacement error, codeType=RS, err=%+v", err)
		}
	case LRC:

	default:

	}
}

func (dcm *DCManager) GeneratePlacementByArchType() error {
	switch dcm.erasureCodeConf.ChunkPlaceType {
	case FLAT:
		if dcm.rackManager.racksNum < dcm.erasureCodeConf.N {
			logrus.Errorf("[DCManager.GenerateRSPlacement] error params for rack init, racksNum=%d,N=%d", dcm.rackManager.racksNum, dcm.erasureCodeConf.N)
			return enum_error.ParamsInvalidError
		}
		for stripeId := 0; stripeId < dcm.stripesNum; stripeId++ {
			rackIdList := util.GenerateListSample(dcm.rackManager.racksNum, dcm.erasureCodeConf.N)
			diskIdList := make([]int, 0)
			for _, rackId := range rackIdList {
				diskId := dcm.GetDiskRandomlyByRack(rackId)
				dcm.disksManager.SetDiskStripe(diskId, stripeId, rackId)
				diskIdList = append(diskIdList, diskId)
			}
			dcm.stripesLocation = append(dcm.stripesLocation, diskIdList)
		}
	case HIERARCHICAL:
	default:

	}
	return nil
}

func (dcm *DCManager) GetDiskRandomlyByRack(rackId int) int {
	minDiskNumber := rackId * dcm.nodesPerRack * dcm.disksPerNode
	maxDiskNumber := minDiskNumber + dcm.nodesPerRack*dcm.disksPerNode - 1
	if minDiskNumber == maxDiskNumber {
		return minDiskNumber
	}
	return util.RandomInt(minDiskNumber, maxDiskNumber)
}
