package internal

import (
	"math/rand"
	"time"
)

type SensorData struct {
	Value  float64   `json:"value"`
	Unit   string    `json:"unit"`
	ID1    int32     `json:"id1"`
	ID2    string    `json:"id2"`
	ReadAt time.Time `json:"readat"`
}

type SensorDataResp struct {
	SuccessfulCnt int64 `json:"successfulcnt"`
}

type Config struct {
	DataChanSize int // The max size of the buffered data channel.
	/*
		The limit of the buffer, the data channel will be flushed once it crosses this limit.
	*/
	BufferLimit   int
	BaseURL       string // The base url to post data to.
	BufferPostURL string // The url to post buffered data to.
	StreamPostURL string // The url to post continuos stream data to.
}

var Configs = Config{
	DataChanSize: 100,
	BufferLimit:  70,
}

func getRandomInRange(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

func getRandomSensorValInRange(min, max float64) float64 {
	return rand.Float64()*(max-min+1) + min
}

var sensorUnits = []string{"Degrees Celsius", "Percentage", "Meters/second squared",
	"Lumen/square meter", "Kilopascal", "Volts", "Amperes", "Parts/million",
	"Decibels", "Grams/cubic meter"}

var sensorValueRanges = map[string][2]float64{
	"Degrees Celsius":       {-40.0, 125.0},
	"Percentage":            {0.0, 100.0},
	"Meters/second squared": {-10.0, 10.0},
	"Lumen/square meter":    {0.0, 100000.0},
	"Kilopascal":            {30.0, 300.0},
	"Volts":                 {0.0, 5.0},
	"Amperes":               {0.0, 50.0},
	"Parts/million":         {0.0, 5000.0},
	"Decibels":              {30.0, 120.0},
	"Grams/cubic meter":     {0.0, 500.0},
}

var subSensorUnitMap = map[string]string{
	"A": "Degrees Celsius",
	"B": "Percentage",
	"C": "Meters/second squared",
	"D": "Lumen/square meter",
	"E": "Kilopascal",
	"F": "Volts",
	"G": "Amperes",
	"H": "Parts/million",
	"I": "Decibels",
	"J": "Grams/cubic meter",
}
