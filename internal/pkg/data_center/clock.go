package data_center

type DeviceState int8

const (
	DeviceStateNormal DeviceState = iota
	DeviceStateUnavailable
	DeviceStateCrashed
)

type DeviceClock struct {
	localTime        float64
	globalTime       float64
	repairTime       float64
	repairStart      float64
	lastUpdateTime   float64
	unavailableStart float64
	unavailableTime  float64
}

func (dc *DeviceClock) Init(currentTime float64) {
	dc.localTime = 0
	dc.globalTime = currentTime
	dc.repairTime = 0
	dc.repairStart = 0
	dc.lastUpdateTime = currentTime
	dc.unavailableStart = 0
	dc.unavailableTime = 0
}

func (dc *DeviceClock) Update(currentTime float64, deviceState DeviceState) {
	dc.localTime += currentTime - dc.lastUpdateTime
	dc.lastUpdateTime = currentTime
	if deviceState == DeviceStateCrashed {
		dc.repairTime = currentTime - dc.repairStart
	} else {
		dc.repairTime = 0
	}
}

func (dc *DeviceClock) GetLocalTime() float64 {
	return dc.localTime
}

func (dc *DeviceClock) GetRepairTime() float64 {
	return dc.repairTime
}
