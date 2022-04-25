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

func (r *Rack) Fail() {
	r.state = RackStateUnavailable
}

func (r *Rack) Repair() {
	r.state = RackStateNormal
}

func (r *Rack) Crash() {
	r.state = RackStateCrashed
}

func (r *Rack) ResetState() {
	r.state = RackStateNormal
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

func (rm *RacksManager) Reset(currentTime float64) {
	for _, rack := range rm.racks {
		rack.ResetState()
	}
}

func (rm *RacksManager) isValidRackId(rackId int) bool {
	return rackId >= 0 && rackId < len(rm.racks)
}

func (rm *RacksManager) GetRackState(rackId int) RackState {
	if rm.isValidRackId(rackId) {
		return rm.racks[rackId].GetState()
	}
	return RackStateUndefined
}

func (rm *RacksManager) FailRack(rackId int) {
	if rm.isValidRackId(rackId) {
		rm.racks[rackId].Fail()
	}
}

func (rm *RacksManager) RepairRack(rackId int) {
	if rm.isValidRackId(rackId) {
		rm.racks[rackId].Repair()
	}
}

func (rm *RacksManager) GetRackRepairDistribution(rackId int) *util.Weibull {
	if rm.isValidRackId(rackId) {
		return rm.racks[rackId].rackRepairDistribution
	}
	return nil
}

func (rm *RacksManager) GetRackFailDistribution(rackId int) *util.Weibull {
	if rm.isValidRackId(rackId) {
		return rm.racks[rackId].rackFailDistribution
	}
	return nil
}

func (rm *RacksManager) GetRackNum() int {
	return len(rm.racks)
}
