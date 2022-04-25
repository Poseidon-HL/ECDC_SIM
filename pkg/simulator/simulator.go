package simulator

import (
	"ECDC_SIM/internal/pkg/data_center"
	"ECDC_SIM/internal/pkg/event_trigger"
	"github.com/gogap/logrus"
)

type Simulator struct {
	dcConf       *data_center.DCConf
	ecConf       *data_center.ErasureCodeConf
	eventManager *event_trigger.EventManager
}

type SimResult struct {
	FailedStripesNum       int
	LostChunkNum           int
	BlockedRatio           float64
	SingleChunkRepairRatio float64
}

func NewSimulator(dcConf *data_center.DCConf, ecConf *data_center.ErasureCodeConf) *Simulator {
	data_center.InitDCManager(dcConf, ecConf)
	return &Simulator{
		dcConf:       dcConf,
		ecConf:       ecConf,
		eventManager: event_trigger.NewEventManager(),
	}
}

func (s *Simulator) Reset() {
	data_center.GetDCManager().Reset()
	s.eventManager.ResetEventManager()
}

func (s *Simulator) RunIteration(iteration int) *SimResult {
	s.Reset()
	var currentTime float64
	logrus.Infof("[Simulator.RunIteration] ite=%d", iteration)
	for {
		eventExecRes := s.eventManager.HandleNextEvent(currentTime)
		if eventExecRes.EventTime > s.dcConf.MissionTime {
			break
		}
		currentTime = eventExecRes.EventTime
		switch eventExecRes.EventType {
		case event_trigger.EventDiskFail, event_trigger.EventNodeFail:
			dataLoss, failedStripesNum, lostChunkNum := data_center.GetDCManager().CheckDataLoss()
			if dataLoss {
				failedStripesNum += s.eventManager.GetDelayedRepairDictLength()
				lostChunkNum += s.eventManager.GetDelayedRepairDictLength()
				return &SimResult{
					FailedStripesNum:       failedStripesNum,
					LostChunkNum:           lostChunkNum,
					BlockedRatio:           data_center.GetDCManager().GetBlockedRatio(currentTime),
					SingleChunkRepairRatio: s.eventManager.GetSingleChunkRepairRatio(),
				}
			}
		}
	}
	logrus.Infof("[Simulator.RunIteration] ite=%d, no data loss happen", iteration)
	return &SimResult{
		BlockedRatio:           data_center.GetDCManager().GetBlockedRatio(currentTime),
		SingleChunkRepairRatio: s.eventManager.GetSingleChunkRepairRatio(),
	}
}
