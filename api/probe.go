package api

import (
	"bytes"
	"compress/gzip"
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

	client *http.Client
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

	p.client = &http.Client{Timeout: 20 * time.Second}

	ticker := p.Clock().NewTicker(config.Config.Cycle)

	for range ticker.Chan() {
		logger.Log.Debugf("Get and send model")

		p.getAndSendModel(session)
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
	var buffer bytes.Buffer
	if config.Config.Gzip {
		g := gzip.NewWriter(&buffer)
		if _, err = g.Write(data); err != nil {
			logger.Log.Debugf("Can't gzip model: %v", err)
			return
		}
		if err = g.Close(); err != nil {
			logger.Log.Debugf("Can't close gzip writer: %v", err)
			return
		}
	} else {
		buffer.Write(data)
	}

	httpRequest, err := http.NewRequest("POST", config.Config.RemoteUrl, &buffer)
	if err != nil {
		logger.Log.Debugf("Error while creating request: %v", err)
		return
	}
	httpRequest.Header.Set("Authorization", fmt.Sprintf("Token token=%v", config.Config.AuthToken))
	httpRequest.Header.Set("Content-Type", "application/x-protobuf")
	httpRequest.Header.Set("Accept", "application/x-protobuf")
	if config.Config.Gzip {
		httpRequest.Header.Set("Content-Encoding", "gzip")
	}
	httpRequest.ContentLength = int64(buffer.Len())

	logger.Log.Debugf("Protobuf body size: %v bytes", buffer.Len())

	// Send http request
	t = time.Now()
	response, err := p.client.Do(httpRequest)
	logger.Log.Debugf("Ara response time: %v", time.Since(t))
	if err != nil {
		logger.Log.Debugf("Error while sending request: %v", err)
		return
	}
	logger.Log.Debugf("Ara response status code: %v", response.Status)
}
