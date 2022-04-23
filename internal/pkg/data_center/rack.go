package data_center

import "ECDC_SIM/internal/pkg/util"

type RackState int8

const (
	RackStateNormal RackState = iota
	RackStateUnavailable
	RackStateCrashed
)

type Rack struct {
	state                  RackState
	rackFailDistribution   *util.Weibull
	rackRepairDistribution *util.Weibull
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
