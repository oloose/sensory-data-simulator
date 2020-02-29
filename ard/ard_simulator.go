package ard

import (
	"math"
	"time"
)

//id = id of a Arduino
//sensorData = Current sensor data values
//refSenData = Desired sensor data values (the Arduino gradually changes the current values towards this)
//stopChan = Channel used to stop sensor data updates
type Arduino struct {
	id         uint8
	sensorData SensorsData
	refSenData SensorsData
	stopChan   chan struct{}
}

//Initializes and returns a new Arduino instance
func NewArduino(mId uint8) *Arduino {
	sim := Arduino{mId, *newSensorData(), SensorsData{}, make(chan struct{})}
	generateRandomValues(&sim.refSenData)

	sim.StartSensorUpdate()
	return &sim
}

//Returns a Arduinos current sensor data values
func (rArd *Arduino) SensorData() SensorsData {
	return rArd.sensorData
}

//Sets a desired sensor data values for a Arduino (every parameter needs to be set)
func (rArd *Arduino) SetRefSensorData(mTemp float64, mHum float64, mPh float64, mWater float64, mLight float64) {
	rArd.refSenData.Temperature = mTemp
	rArd.refSenData.Humidity = mHum
	rArd.refSenData.Ph = mPh
	rArd.refSenData.Waterlevel = mWater
	rArd.refSenData.Lightintensity = mLight
}

//Starts the gradual update of sensor data towards the desired sensor data values
func (rArd *Arduino) StartSensorUpdate() {
	pauseChan := make(chan bool)
	rArd.sensorData.graduallyChangeAll(rArd.stopChan, pauseChan, &rArd.refSenData)

	go func() {
		select {
		default:
			for {
				time.Sleep(1 * time.Second)
				desValCount := 0

				//check if sensor values have reached desired values
				if math.Abs(rArd.SensorData().Temperature-rArd.refSenData.Temperature) <= tempStep {
					desValCount++
				}

				if math.Abs(rArd.SensorData().Humidity-rArd.refSenData.Humidity) <= humStep {
					desValCount++
				}

				if math.Abs(rArd.SensorData().Ph-rArd.refSenData.Ph) <= phStep {
					desValCount++
				}

				if math.Abs(rArd.SensorData().Waterlevel-rArd.refSenData.Waterlevel) <= waterStep {
					desValCount++
				}

				if math.Abs(rArd.SensorData().Lightintensity-rArd.refSenData.Lightintensity) <= lightStep {
					desValCount++
				}

				// if at least 3 sensor values have reached its desired values, wait for 10s,
				// then set new random desired values
				if desValCount >= 3 {
					pauseChan <- true
					desValCount = 0 // reset counter (needs to be before time.Sleep to prevent deadlock)

					time.Sleep(10 * time.Second) // keep desired sensor data state for 10s, than set new desired values

					generateRandomValues(&rArd.refSenData) //set new desired sensor data
					pauseChan <- false
				}
			}
		case <-rArd.stopChan:
			return
		}
	}()

}

//Stops the gradual updates of sensor data values
func (rArd *Arduino) StopSensorUpdate() {
	close(rArd.stopChan)                // close channel to notify receiver select in graduallyChangeAll/Single() func
	rArd.stopChan = make(chan struct{}) // re-initialize stop chan to prepare new StartSensorUpdate() func call
}

//Returns the id
func (rArd *Arduino) Id() uint8 {
	return rArd.id
}
