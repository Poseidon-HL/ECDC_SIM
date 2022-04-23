package data_center

import (
	"ECDC_SIM/internal/pkg/util"
)

type DiskState int8

const (
	DiskStateNormal DiskState = iota
	DiskStateUnavailable
	DiskStateCrashed
)

type Disk struct {
	stripeId               []int
	stripeIndex            []int
	state                  DiskState
	diskFailDistribution   *util.Weibull
	diskRepairDistribution *util.Weibull
}

type DisksManager struct {
	disksNum           int
	disks              []*Disk
	failedDiskNum      int
	unavailableDiskNum int
	failedDiskMap      map[int]int
	unavailableDiskMap map[int]int
	diskCapacity       int
}

func NewDisksManager(disksNum, diskCap int, dFailD, dRepairD *util.Weibull) *DisksManager {
	disksManager := &DisksManager{
		disksNum:           disksNum,
		failedDiskMap:      make(map[int]int),
		unavailableDiskMap: make(map[int]int),
		diskCapacity:       diskCap,
	}
	for i := 0; i < disksNum; i++ {
		disksManager.disks = append(disksManager.disks, &Disk{
			state:                  DiskStateNormal,
			diskFailDistribution:   dFailD,
			diskRepairDistribution: dRepairD,
		})
	}
	return disksManager
}

func (dm *DisksManager) SetDiskStripe(diskId, stripeId, stripeIdx int) {
	dm.disks[diskId].stripeId, dm.disks[diskId].stripeIndex = append(dm.disks[diskId].stripeId, stripeId), append(dm.disks[diskId].stripeIndex, stripeIdx)
}
