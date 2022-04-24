package data_center

import (
	"ECDC_SIM/internal/pkg/util"
)

type DiskState int8

const (
	DiskStateNormal DiskState = iota
	DiskStateUnavailable
	DiskStateCrashed
	Undefined
)

type Disk struct {
	diskClock              *DeviceClock
	stripeId               []int
	stripeIndex            []int
	state                  DiskState
	diskFailDistribution   *util.Weibull
	diskRepairDistribution *util.Weibull
}

func (d *Disk) ResetState() {
	d.state = DiskStateNormal
}

func (d *Disk) GetState() DiskState {
	return d.state
}

func (d *Disk) Fail(currentTime float64) {
	if d.state == DiskStateNormal {
		d.diskClock.unavailableStart = currentTime
	}
	d.state = DiskStateCrashed
	d.diskClock.repairStart = currentTime
	d.diskClock.repairTime = 0
}

func (d *Disk) Offline(currentTime float64) {
	if d.state == DiskStateNormal {
		d.state = DiskStateUnavailable
		d.diskClock.unavailableStart = currentTime
	}
}

func (d *Disk) Online(currentTime float64) {
	if d.state == DiskStateUnavailable {
		d.state = DiskStateNormal
		d.diskClock.unavailableTime += currentTime - d.diskClock.unavailableStart
	}
}

func (d *Disk) GetStripes() []int {
	return d.stripeId
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
			diskClock:              new(DeviceClock),
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

func (dm *DisksManager) Reset(currentTime float64) {
	for _, disk := range dm.disks {
		disk.diskClock.Init(currentTime)
		disk.ResetState()
	}
}

func (dm *DisksManager) isValidDiskId(diskId int) bool {
	return diskId >= 0 && diskId < len(dm.disks)
}

func (dm *DisksManager) GetDiskState(diskId int) DiskState {
	if dm.isValidDiskId(diskId) {
		return dm.disks[diskId].GetState()
	}
	return Undefined
}

func (dm *DisksManager) FailDisk(diskId int, currentTime float64) {
	if dm.isValidDiskId(diskId) {
		dm.disks[diskId].Fail(currentTime)
	}
}

func (dm *DisksManager) OfflineDisk(diskId int, currentTime float64) {
	if dm.isValidDiskId(diskId) {
		dm.disks[diskId].Offline(currentTime)
	}
}

func (dm *DisksManager) OnlineDisk(diskId int, currentTime float64) {
	if dm.isValidDiskId(diskId) {
		dm.disks[diskId].Online(currentTime)
	}
}

func (dm *DisksManager) GetDiskStripes(diskId int) []int {
	if dm.isValidDiskId(diskId) {
		return dm.disks[diskId].GetStripes()
	}
	return nil
}
