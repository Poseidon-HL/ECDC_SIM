package simulator

import (
	"ECDC_SIM/internal/pkg/data_center"
	"ECDC_SIM/internal/pkg/event_trigger"
)

type Simulator struct {
	eventManager *event_trigger.EventManager
}

func (s *Simulator) MustInit(dcConf *data_center.DCConf, ecConf *data_center.ErasureCodeConf) {
	data_center.InitDCManager(dcConf, ecConf)
	s.eventManager = event_trigger.NewEventManager()
}
