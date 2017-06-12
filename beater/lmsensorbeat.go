package beater

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/mdlayher/lmsensors"
	"github.com/singlehopllc/lmsensorbeat/config"
)

type Lmsensorbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Lmsensorbeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Lmsensorbeat) getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func (bt *Lmsensorbeat) Run(b *beat.Beat) error {
	logp.Info("lmsensorbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	scanner := lmsensors.New()
	devices, err := scanner.Scan()
	if err != nil {
		logp.Err("Error Scanning For Devices: %s", err)
		return err
	}
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
		if len(devices) == 0 {
			return errors.New("No Devices Found!")
		}
		var events []common.MapStr
		//var event common.MapStr
		for _, device := range devices {
			deviceName := device.Name
			for _, sensor := range device.Sensors {
				sdata := map[string]interface{}{
					"name":   deviceName,
					"sensor": sensor,
				}
				stype := bt.getType(sensor)
				logp.Info("Sensor Type: %s", bt.getType(sensor))
				event := common.MapStr{
					"@timestamp": common.Time(time.Now()),
					"type":       stype,
					"device":     sdata,
				}
				events = append(events, event)
			}
		}
		bt.client.PublishEvents(events)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Lmsensorbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
