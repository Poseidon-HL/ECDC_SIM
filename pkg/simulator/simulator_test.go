package simulator

import (
	"ECDC_SIM/internal/pkg/data_center"
	"ECDC_SIM/internal/pkg/event_trigger"
	"ECDC_SIM/internal/pkg/util"
	"github.com/gogap/logrus"
	"math"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestNewSimulatorRS(t *testing.T) {
	type args struct {
		dcConf *data_center.DCConf
		ecConf *data_center.ErasureCodeConf
		rConf  *event_trigger.RunningConfig
	}
	tests := []struct {
		name string
		args args
		want *Simulator
	}{
		{
			name: "TestSimulatorRS", args: args{
				dcConf: &data_center.DCConf{
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
				ecConf: &data_center.ErasureCodeConf{
					CodeType:       data_center.RS,
					ChunkPlaceType: data_center.FLAT,
					N:              9,
					K:              6,
				},
				rConf: &event_trigger.RunningConfig{
					UseTrace:               false,
					UsePowerOutage:         false,
					EnableTransientFailure: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logFile, _ := os.OpenFile("../../output/test_log/"+strconv.Itoa(int(time.Now().UnixNano()))+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
			logrus.SetOutput(logFile)
			got := NewSimulator(tt.args.dcConf, tt.args.ecConf, tt.args.rConf)
			iteration := 1
			for ite := 0; ite < iteration; ite++ {
				result := got.RunIteration(ite)
				t.Log(result)
			}
		})
	}
}

func TestNewSimulatorLRC(t *testing.T) {
	type args struct {
		dcConf *data_center.DCConf
		ecConf *data_center.ErasureCodeConf
		rConf  *event_trigger.RunningConfig
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "TestSimulatorLRC", args: args{
			dcConf: &data_center.DCConf{
				RacksNum:                    32,
				StripesNum:                  340000,
				DisksPerNode:                1,
				DiskCapacity:                int(math.Pow(2, 10)),
				NodesPerRack:                72,
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
			ecConf: &data_center.ErasureCodeConf{
				CodeType:             data_center.LRC,
				ChunkPlaceType:       data_center.FLAT,
				N:                    16,
				K:                    12,
				L:                    2,
				LRCDataChunkOffset:   [][]int{{0, 1, 2, 3, 4, 5}, {8, 9, 10, 11, 12, 13}},
				LRCLocalChunkParity:  []int{6, 14},
				LRCGlobalChunkParity: []int{7, 15},
			},
			rConf: &event_trigger.RunningConfig{
				UseTrace:               false,
				UsePowerOutage:         false,
				EnableTransientFailure: false,
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logFile, _ := os.OpenFile("../../output/test_log/"+strconv.Itoa(int(time.Now().UnixNano()))+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
			logrus.SetOutput(logFile)
			got := NewSimulator(tt.args.dcConf, tt.args.ecConf, tt.args.rConf)
			iteration := 1
			for ite := 0; ite < iteration; ite++ {
				result := got.RunIteration(ite)
				t.Log(result)
			}
		})
	}
}
