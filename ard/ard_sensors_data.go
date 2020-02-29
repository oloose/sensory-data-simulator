package ard

import (
	"math/rand"
	"time"
)

// Min/Max sensor values
const tempMin = 10
const tempMax = 40
const humMin = 20
const humMax = 100
const phMin = 0
const phMax = 14
const waterMin = 0
const waterMax = 10
const lightMin = 0
const lightMax = 100

// Change Step values
const tempStep = 0.25
const humStep = 2
const phStep = 0.05
const waterStep = 0.25
const lightStep = 1

type SensorsData struct {
	Temperature    float64 // in celsius e.g. 23°C
	Humidity       float64 // in percent e.g. 90% Humidity
	Ph             float64 // as 5.6 format
	Waterlevel     float64 // in centimeters of water height in basin (from ground to water surface)
	Lightintensity float64 // percent from 100 (max, all lamps working) to 0 (min, all lamps shattered)
}

//Creates new sensor data an initializes it with random values
func newSensorData() *SensorsData {
	sensorsData := SensorsData{}
	generateRandomValues(&sensorsData)
	return &sensorsData
}

//mSenData = sensor data strcut the random values will be insert in
func generateRandomValues(mSenData *SensorsData) {
	mSenData.Temperature = tempMin + rand.Float64()*(tempMax-tempMin)
	mSenData.Humidity = humMin + rand.Float64()*(humMax-humMin)
	mSenData.Ph = phMin + rand.Float64()*(phMax-phMin)
	mSenData.Waterlevel = waterMin + rand.Float64()*(waterMax-waterMin)
	mSenData.Lightintensity = lightMin + rand.Float64()*(lightMax-lightMin)
}

//Starts go routines for each value (Temperatur, Humidity, Ph, Waterlevel,
// Lightintensity) dsa handles the gradual change
//mStopChan = channel to stop the gradual change
//mPausChan = channel to pause the gradual change (true = pause, false = run)
//mSenRef = Reference sensor data with desired sensor values
func (rSenData *SensorsData) graduallyChangeAll(mStopChan <-chan struct{}, mPauseChan <-chan bool,
	mSenRef *SensorsData) {
	go rSenData.graduallyChangeSingle(mStopChan, mPauseChan, &rSenData.Temperature, &mSenRef.Temperature, tempMin, tempMax,
		tempStep)
	go rSenData.graduallyChangeSingle(mStopChan, mPauseChan, &rSenData.Humidity, &mSenRef.Humidity, humMin, humMax, humStep)
	go rSenData.graduallyChangeSingle(mStopChan, mPauseChan, &rSenData.Ph, &mSenRef.Ph, phMin, phMax, phStep)
	go rSenData.graduallyChangeSingle(mStopChan, mPauseChan, &rSenData.Waterlevel, &mSenRef.Waterlevel, waterMin, waterMax,
		waterStep)
	go rSenData.graduallyChangeSingle(mStopChan, mPauseChan, &rSenData.Lightintensity, &mSenRef.Lightintensity, lightMin,
		lightMax, lightStep)
}

//mStopChan = channel to receive stop
//mCur = Current sensor value
//mNew = desired sensor value
//mMin = minimum sensor value (e.g. not colder than 10°C)
//mMax = maximum sensor value  (e.g. not hotter than 40°C)
//mStep = Factor to decrease/increase senosr value by
func (rSenData *SensorsData) graduallyChangeSingle(mStopChan <-chan struct{}, mPauseChan <-chan bool, mCur *float64,
	mNew *float64, mMin float64, mMax float64, mStep float64) {
	b := true
	pause := false

	//change value gradually every second
	for {
		time.Sleep(500 * time.Millisecond)

		select {
		default:
			if !pause {
				//decide to gradually lower or higher value
				if *mNew > *mCur {
					b = true
				} else {
					b = false
				}
				if b {
					if (*mCur+mStep < mMax) && (*mCur < *mNew) {
						*mCur = *mCur + mStep
					}
				} else {
					if (*mCur-mStep > mMin) && (*mCur > *mNew) {
						*mCur = *mCur - mStep
					}
				}
			}

		case <-mStopChan:
			return

		case p := <-mPauseChan:
			pause = p
		}
	}
}
