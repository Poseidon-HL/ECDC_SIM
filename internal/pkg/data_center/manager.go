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
	powerOutageD    *util.Weibull
}

type DCConf struct {
	RacksNum                    int
	StripesNum                  int
	DisksPerNode                int
	DiskCapacity                int
	NodesPerRack                int
	ChunkNum                    int
	ChunkSize                   int
	DataChunksNum               int
	NFailD, NTFailD, NTRepairD  *util.Weibull
	DFailD, DRepairD            *util.Weibull
	RFailD, RRepairD            *util.Weibull
	powerOutageD                *util.Weibull
	MaxCrossRackRepairBandwidth float64
	MaxIntraRackRepairBandwidth float64
	MissionTime                 float64
	UseNetwork                  bool
}

func InitDCManager(dcConf *DCConf, eCConf *ErasureCodeConf) {
	dcManager = &DCManager{
		state:           OK,
		disksPerNode:    dcConf.DisksPerNode,
		nodesPerRack:    dcConf.NodesPerRack,
		stripesNum:      dcConf.StripesNum,
		chunksNum:       dcConf.ChunkNum,
		chunkSize:       dcConf.ChunkSize,
		dataChunksNum:   dcConf.DataChunksNum,
		erasureCodeConf: eCConf,
		missionTime:     dcConf.MissionTime,
	}
	dcManager.nodesManager = NewNodesManager(dcConf.RacksNum*dcConf.NodesPerRack, dcConf.NFailD, dcConf.NTFailD, dcConf.NTRepairD)
	dcManager.disksManager = NewDisksManager(dcManager.nodesManager.nodesNum*dcConf.DisksPerNode, dcConf.DiskCapacity, dcConf.DFailD, dcConf.DRepairD)
	dcManager.rackManager = NewRacksManager(dcConf.RacksNum, dcConf.RFailD, dcConf.RRepairD)
	dcManager.networkManager = NewNetworkManager(dcConf.RacksNum, dcConf.UseNetwork, dcConf.MaxCrossRackRepairBandwidth, dcConf.MaxIntraRackRepairBandwidth)
	return
}

func GetDCManager() *DCManager {
	return dcManager
}

func (dcm *DCManager) Reset() {
	dcm.disksManager.Reset(0)
	dcm.nodesManager.Reset(0)
	dcm.rackManager.Reset(0)
	dcm.networkManager.Reset()
	dcm.GenerateDataPlacement()
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

func (dcm *DCManager) GetNodeIdByRackId(rackId int, offset int) int {
	return rackId*dcm.nodesPerRack + offset
}

func (dcm *DCManager) GetNodeIdByDiskId(diskId int) int {
	return diskId / dcm.disksPerNode
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

func (dcm *DCManager) GetPowerOutageD() *util.Weibull {
	return dcm.powerOutageD
}

// GenerateDataPlacement 生成数据块放置策略
func (dcm *DCManager) GenerateDataPlacement() {
	var err error
	switch dcm.erasureCodeConf.CodeType {
	case RS, LRC:
		logrus.Info("[DCManager.GenerateDataPlacement] generate placement for code RS")
		err = dcm.GeneratePlacementByArchType()
		if err != nil {
			logrus.Errorf("DCManager.GenerateDataPlacement error, codeType=RS, err=%+v", err)
		}
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
		logrus.Error("[DCManager.GeneratePlacementByArchType] invalid chunk place type")
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

func (dcm *DCManager) CheckDataLoss() (bool, int, int) {
	failedDiskMap := dcm.disksManager.GetFailedDiskMap()
	// TODO check logic here
	stripeIdList := make([]int, 0)
	for _, failedDisk := range failedDiskMap {
		stripeIdList = append(stripeIdList, dcm.disksManager.GetDiskStripes(failedDisk)...)
	}
	var dataLoss bool
	var failedStripes int
	var lostChunks int
	switch dcm.erasureCodeConf.CodeType {
	case RS:
		for _, stripeId := range stripeIdList {
			curStripeFailedDiskNum := 0
			curStripeLostChunksNum := 0
			for _, stripeDiskId := range dcm.stripesLocation[stripeId] {
				if _, ok := failedDiskMap[stripeDiskId]; ok {
					curStripeFailedDiskNum += 1
					curStripeLostChunksNum += 1
				}
			}
			if curStripeFailedDiskNum > dcm.erasureCodeConf.N-dcm.erasureCodeConf.K {
				dataLoss = true
				failedStripes += 1
				lostChunks += curStripeLostChunksNum
			}
		}
		return dataLoss, failedStripes, lostChunks
	case LRC:
	}

	return false, 0, 0
}

func (dcm *DCManager) GetBlockedRatio(currentTime float64) float64 {
	sumOfUnavailingTime := dcm.disksManager.GetSumOfDiskUnavailableTime(currentTime)
	logrus.Infof("[GetBlockedRatio] sumOfUnavailingTime=%+v", sumOfUnavailingTime)
	return sumOfUnavailingTime / (float64(dcm.chunksNum) * currentTime)
}
