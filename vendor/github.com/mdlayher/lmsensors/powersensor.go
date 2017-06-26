package lmsensors

import (
	"strconv"
	"time"
)

var _ Sensor = &PowerSensor{}

// A PowerSensor is a Sensor that detects average electrical power consumption
// in watts.
type PowerSensor struct {
	// The name of the sensor.
	Name string `json:"name"`

	// The average electrical power consumption, in watts, indicated
	// by the sensor.
	Average float64 `json:"average"`

	// The interval of time over which the average electrical power consumption
	// is collected.
	AverageInterval time.Duration `json:"average_interval"`

	// Whether or not this sensor has a battery.
	Battery bool `json:"battery"`

	// The model number of the sensor.
	ModelNumber string `json:"model_number"`

	// Miscellaneous OEM information about the sensor.
	OEMInfo string `json:"oem_info"`

	// The serial number of the sensor.
	SerialNumber string `json:"serial_number"`
}

func (s *PowerSensor) name() string        { return s.Name }
func (s *PowerSensor) setName(name string) { s.Name = name }

func (s *PowerSensor) parse(raw map[string]string) error {
	for k, v := range raw {
		switch k {
		case "average":
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}

			// Raw temperature values are scaled by one million
			f /= 1000000
			s.Average = f
		case "average_interval":
			// Time values in milliseconds
			d, err := time.ParseDuration(v + "ms")
			if err != nil {
				return err
			}

			s.AverageInterval = d
		case "is_battery":
			s.Battery = v != "0"
		case "model_number":
			s.ModelNumber = v
		case "oem_info":
			s.OEMInfo = v
		case "serial_number":
			s.SerialNumber = v
		}
	}

	return nil
}
