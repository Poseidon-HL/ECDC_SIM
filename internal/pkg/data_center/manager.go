package data_center

import (
	"ECDC_SIM/internal/pkg/enum_error"
	"ECDC_SIM/internal/pkg/util"
	"github.com/gogap/logrus"
)

var (
	dcManager *DCManager
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
	networkManager  *NetworkManager
	disksPerNode    int // 每一节点上的磁盘数
	nodesPerRack    int // 每一机架上的节点数
	stripesNum      int
	chunksNum       int
	chunkSize       int
	dataChunksNum   int
	erasureCodeConf *ErasureCodeConf
	stripesLocation [][]int
	missionTime     float64
}

type DCConf struct {
	racksNum                    int
	stripesNum                  int
	disksPerNode                int
	diskCapacity                int
	nodesPerRack                int
	chunkNum                    int
	chunkSize                   int
	dataChunksNum               int
	nFailD, nTFailD, nTRepairD  *util.Weibull
	dFailD, dRepairD            *util.Weibull
	rFailD, rRepairD            *util.Weibull
	maxCrossRackRepairBandwidth float64
	maxIntraRackRepairBandwidth float64
	missionTime                 float64
}

func InitDCManager(dcConf *DCConf, eCConf *ErasureCodeConf) {
	dcManager = &DCManager{
		state:           OK,
		disksPerNode:    dcConf.disksPerNode,
		nodesPerRack:    dcConf.nodesPerRack,
		stripesNum:      dcConf.stripesNum,
		chunksNum:       dcConf.chunkNum,
		chunkSize:       dcConf.chunkSize,
		dataChunksNum:   dcConf.dataChunksNum,
		erasureCodeConf: eCConf,
		missionTime:     dcConf.missionTime,
	}
	dcManager.nodesManager = NewNodesManager(dcConf.racksNum*dcConf.nodesPerRack, dcConf.nFailD, dcConf.nTFailD, dcConf.nTRepairD)
	dcManager.disksManager = NewDisksManager(dcManager.nodesManager.nodesNum*dcConf.disksPerNode, dcConf.diskCapacity, dcConf.dFailD, dcConf.dRepairD)
	dcManager.rackManager = NewRacksManager(dcConf.racksNum, dcConf.rFailD, dcConf.rRepairD)
	dcManager.networkManager = NewNetworkManager(dcConf.racksNum, dcConf.maxCrossRackRepairBandwidth, dcConf.maxIntraRackRepairBandwidth)
	return
}

func GetDCManager() *DCManager {
	return dcManager
}

func (dcm *DCManager) DiskManager() *DisksManager {
	return dcm.disksManager
}

func (dcm *DCManager) NodeManager() *NodesManager {
	return dcm.nodesManager
}

func (dcm *DCManager) RackManager() *RacksManager {
	return dcm.rackManager
}

func (dcm *DCManager) Network() *NetworkManager {
	return dcm.networkManager
}

func (dcm *DCManager) ErasureCodeConf() *ErasureCodeConf {
	return dcm.erasureCodeConf
}

func (dcm *DCManager) GetRackIdByDiskId(diskId int) int {
	return diskId / dcm.nodesPerRack * dcm.disksPerNode
}

func (dcm *DCManager) GetDiskIdByNodeId(nodeId int, offset int) int {
	return nodeId*dcm.disksPerNode + offset
}

func (dcm *DCManager) isValidStripeId(stripeId int) bool {
	return stripeId >= 0 && stripeId < len(dcm.stripesLocation)
}

func (dcm *DCManager) GetStripesLocation(stripeId int) []int {
	if dcm.isValidStripeId(stripeId) {
		return dcm.stripesLocation[stripeId]
	}
	return nil
}

func (dcm *DCManager) GetNodesPerRack() int {
	return dcm.nodesPerRack
}

func (dcm *DCManager) GetDisksPerNode() int {
	return dcm.disksPerNode
}

func (dcm *DCManager) GetChunkSize() int {
	return dcm.chunkSize
}

func (dcm *DCManager) GetMissionTime() float64 {
	return dcm.missionTime
}

// GenerateDataPlacement 生成数据块放置策略
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
