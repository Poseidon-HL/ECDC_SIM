package data_center

import "ECDC_SIM/internal/pkg/util"

type RackState int8

const (
	RackStateNormal RackState = iota
	RackStateUnavailable
	RackStateCrashed
	RackStateUndefined
)

type Rack struct {
	rackClock              *DeviceClock
	state                  RackState
	rackFailDistribution   *util.Weibull
	rackRepairDistribution *util.Weibull
}

func (r *Rack) GetState() RackState {
	return r.state
}

type RacksManager struct {
	racksNum       int
	racks          []*Rack
	failedRacksNum int
}

func NewRacksManager(racksNum int, rFailD, rRepairD *util.Weibull) *RacksManager {
	racksManager := &RacksManager{
		racksNum:       racksNum,
		failedRacksNum: 0,
	}
	for i := 0; i < racksNum; i++ {
		racksManager.racks = append(racksManager.racks, &Rack{
			state:                  RackStateNormal,
			rackFailDistribution:   rFailD,
			rackRepairDistribution: rRepairD,
		})
	}
	return racksManager
}

func (rm *RacksManager) isValidDiskId(rackId int) bool {
	return rackId >= 0 && rackId < len(rm.racks)
}

func (rm *RacksManager) GetRackState(rackId int) RackState {
	if rm.isValidDiskId(rackId) {
		return rm.racks[rackId].GetState()
	}
	return RackStateUndefined
}
