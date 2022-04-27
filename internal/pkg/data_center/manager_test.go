package data_center

import (
	"ECDC_SIM/internal/pkg/util"
	"math"
	"testing"
)

func TestInitDCManager(t *testing.T) {
	type args struct {
		dcConf *DCConf
		eCConf *ErasureCodeConf
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "TestGetStripeLocation", args: args{
			dcConf: &DCConf{
				RacksNum:                    32,
				StripesNum:                  340000,
				DisksPerNode:                1,
				DiskCapacity:                int(math.Pow(2, 10)),
				NodesPerRack:                32,
				ChunkNum:                    340000 * 9,
				ChunkSize:                   256,
				DataChunksNum:               0,
				NFailD:                      util.NewWeibull(1, 91250, 0),
				NTFailD:                     util.NewWeibull(1, 2890.8, 0),
				NTRepairD:                   util.NewWeibull(1, 0.25, 0),
				DFailD:                      util.NewWeibull(1.12, 87600, 0),
				DRepairD:                    nil,
				RFailD:                      util.NewWeibull(1.0, 87600, 0),
				RRepairD:                    util.NewWeibull(1.0, 24, 10),
				MaxCrossRackRepairBandwidth: 125,
				MaxIntraRackRepairBandwidth: 125,
				MissionTime:                 87600,
				UseNetwork:                  true,
			},
			eCConf: &ErasureCodeConf{
				CodeType:       RS,
				ChunkPlaceType: FLAT,
				N:              9,
				K:              6,
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitDCManager(tt.args.dcConf, tt.args.eCConf)
			dcm := GetDCManager()
			dcm.Reset()
			stripeId := 1
			for _, diskId := range dcm.GetStripesLocation(stripeId) {
				for _, stripe := range dcm.disksManager.GetDiskStripes(diskId) {
					if stripe == stripeId {
						t.Logf("disk num =%d, stripeID=%d", diskId, stripe)
					}
				}
			}
		})
	}
}
