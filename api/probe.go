package api

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	external_models "github.com/af83/ara-external-models"
	"github.com/af83/scops/clock"
	"github.com/af83/scops/config"
	"github.com/af83/scops/logger"
	"github.com/gocraft/dbr"
	"github.com/golang/protobuf/proto"
)

type Feeder interface {
	DbConnect() *dbr.Session
	GetCompleteModel(sess *dbr.Session) (*external_models.ExternalCompleteModel, error)
}

type Probe struct {
	clock.ClockConsumer

	feeder Feeder
}

func NewProbe(feeder Feeder) *Probe {
	return &Probe{feeder: feeder}
}

func (p Probe) Run() {
	logger.Log.Printf("Starting the probe")
	logger.Log.Printf("Connecting to the database")
	session := p.feeder.DbConnect()
	logger.Log.Printf("Connection successful")

	tick := p.Clock().After(config.Config.Cycle)

	for {
		select {
		case <-tick:
			logger.Log.Debugf("Get and send model")

			p.getAndSendModel(session)

			tick = p.Clock().After(config.Config.Cycle)
		}
	}
}

func (p Probe) getAndSendModel(s *dbr.Session) {
	t := time.Now()
	model, err := p.feeder.GetCompleteModel(s)
	logger.Log.Debugf("Plugin response time: %v", time.Since(t))
	if err != nil {
		logger.Log.Debugf("Error while fetching model: %v", err)
		return
	}

	logger.Log.Debugf("StopAreas returned: %v", len(model.StopAreas))
	logger.Log.Debugf("Lines returned: %v", len(model.Lines))
	logger.Log.Debugf("VehicleJourneys returned: %v", len(model.VehicleJourneys))
	logger.Log.Debugf("StopVisits returned: %v", len(model.StopVisits))

	data, err := proto.Marshal(model)
	if err != nil {
		logger.Log.Debugf("Error while marshaling model: %v", err)
		return
	}

	// Create http request
	buffer := bytes.NewBuffer(data)
	httpRequest, err := http.NewRequest("POST", config.Config.RemoteUrl, buffer)
	if err != nil {
		logger.Log.Debugf("Error while creating request: %v", err)
		return
	}
	httpRequest.Header.Set("Authorization", fmt.Sprintf("Token token=%v", config.Config.AuthToken))
	httpRequest.Header.Set("Content-Type", "application/x-protobuf")
	httpRequest.Header.Set("Accept", "application/x-protobuf")
	httpRequest.ContentLength = int64(buffer.Len())

	logger.Log.Debugf("Protobuf body size: %v bytes", buffer.Len())

	// Send http request
	httpClient := &http.Client{Timeout: 5 * time.Second}
	t = time.Now()
	response, err := httpClient.Do(httpRequest)
	logger.Log.Debugf("Ara response time: %v", time.Since(t))
	if err != nil {
		logger.Log.Debugf("Error while sending request: %v", err)
		return
	}
	logger.Log.Debugf("Ara response status code: %v", response.Status)
}
