package internal

import (
	"fmt"
	"time"
)


func (g *Generator) SpawnSensor(id1 int32) {
	defer g.wg.Done()

	totalGenCnt := int64(0)

	sleepFor := time.Duration(getRandomInRange(10, 3000)) * time.Millisecond
	ticker := time.NewTicker(sleepFor)

	id2Range := 65 + getRandomInRange(0, 3)
	currID2 := 64
	id2 := ""

	unit := sensorUnits[getRandomInRange(0, 9)]
	minVal, maxVal := g.getMaxMinVal(unit)

	for {

		select {
		case <-g.ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:

			currID2++
			if currID2 > int(id2Range) {
				currID2 = 65
			}
			id2 = string(rune(currID2))

			// This might seem to cause problems but it most probably will never,
			// if the buffer flush limit is kept relatively lower so as to accomodate any late comers.
			dataChnPtr := g.dataChan.Load()

			*dataChnPtr <- &SensorData{
				Value:  getRandomSensorValInRange(minVal, maxVal),
				Unit:   subSensorUnitMap[id2],
				ID1:    id1,
				ID2:    id2,
				ReadAt: time.Now().UTC(),
			}

			totalGenCnt++

			fmt.Printf("\033[%d;0H\033[2KActive : Sensor ID : %7d : Subsensor ID : %7s : Generated : %7d : Sleeps for (ms) : %10d",
				id1+6,
				id1,
				id2,
				totalGenCnt,
				sleepFor.Milliseconds())
		}
	}
}

func (g *Generator) getMaxMinVal(unit string) (minVal, maxVal float64) {

	minVal = sensorValueRanges[unit][0]
	maxVal = sensorValueRanges[unit][1]

	initVal := getRandomSensorValInRange(minVal, maxVal)
	variable := (maxVal - minVal) * 0.10 // +-10%, values oscillate in a narrow band, more realistic.

	if initVal-variable >= minVal {
		minVal = initVal - variable
	}
	if initVal+variable <= maxVal {
		maxVal = initVal + variable
	}
	return
}
