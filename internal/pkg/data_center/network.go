package data_center

type NetworkManager struct {
	useNetwork bool

	maxCrossRackRepairBandwidth   float64
	maxIntraRackRepairBandwidth   float64
	availCrossRackRepairBandwidth float64
	availIntraRackRepairBandwidth []float64
}

func NewNetworkManager(numOfRacks int, useNetwork bool, maxCrossRackRepairBandwidth, maxIntraRackRepairBandwidth float64) *NetworkManager {
	network := &NetworkManager{
		useNetwork:                    useNetwork,
		maxCrossRackRepairBandwidth:   maxCrossRackRepairBandwidth,
		maxIntraRackRepairBandwidth:   maxIntraRackRepairBandwidth,
		availCrossRackRepairBandwidth: maxCrossRackRepairBandwidth,
	}
	for i := 0; i < int(numOfRacks); i++ {
		network.availIntraRackRepairBandwidth = append(network.availIntraRackRepairBandwidth, maxIntraRackRepairBandwidth)
	}
	return network
}

func (n *NetworkManager) Reset() {
	n.availCrossRackRepairBandwidth = n.maxCrossRackRepairBandwidth
	for idx := range n.availIntraRackRepairBandwidth {
		n.availIntraRackRepairBandwidth[idx] = n.maxIntraRackRepairBandwidth
	}
}

func (n *NetworkManager) UpdateAvailCrossRackRepairBandwidth(newBandwidth float64) {
	if newBandwidth <= n.maxCrossRackRepairBandwidth {
		n.availCrossRackRepairBandwidth = newBandwidth
	}
}

func (n *NetworkManager) GetAvailCrossRackRepairBandwidth() float64 {
	return n.availCrossRackRepairBandwidth
}

func (n *NetworkManager) GetAvailIntraRackRepairBandwidth(rackId int) float64 {
	return n.availIntraRackRepairBandwidth[rackId]
}

func (n *NetworkManager) UseNetwork() bool {
	return n.useNetwork
}
