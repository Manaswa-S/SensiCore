package internal

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

type Generator struct {
	/*
		The default channel to send data to.
	*/
	dataChan atomic.Pointer[chan *SensorData]

	/*
		The alternate channel used for flush buffering.
		Main data is to always be pushed to dataChan.
	*/
	altDataChan atomic.Pointer[chan *SensorData]

	/*
		Set to true if data is to be streamed, default is batched.
	*/
	StreamFlush bool

	/*
		The exact number of sensors to be spawned.
		Default is any random number between 5 and 15.
	*/
	SensorCount int32

	/*
		This is the latest ID1 that HAS BEEN USED.
		Always add 1 to it before using.
		By default, ID1 is supposed to be in serial from 1 to {SensorCount}, but can always be overridden.
	*/
	LatestID1 atomic.Int32

	// The same cancel context is used for all go routines.
	ctx       context.Context
	ctxCancel context.CancelCauseFunc
	// The wait group to manage all go routines.
	wg sync.WaitGroup

	/*
		The base config for the generator to work.
		Some fields are required.
	*/
	Configs Config

	/*
		The main channel to push errors to.
	*/
	errChan chan error
}

func NewSensorGenerator(inp *Generator) (*Generator, error) {

	gen := new(Generator)

	dChn := make(chan *SensorData, Configs.DataChanSize)
	gen.dataChan.Store(&dChn)
	adChn := make(chan *SensorData, Configs.DataChanSize)
	gen.altDataChan.Store(&adChn)
	gen.errChan = make(chan error, 10)

	gen.ctx, gen.ctxCancel = context.WithCancelCause(context.Background())

	if inp != nil {
		gen.StreamFlush = inp.StreamFlush

		if inp.SensorCount == 0 {
			gen.SensorCount = int32(getRandomInRange(5, 15))
		} else {
			gen.SensorCount = inp.SensorCount
		}

		gen.LatestID1.Store(inp.LatestID1.Load())

		if inp.Configs.DataChanSize == 0 {
			gen.Configs.DataChanSize = Configs.DataChanSize
		}

		if inp.Configs.BufferLimit == 0 {
			gen.Configs.BufferLimit = Configs.BufferLimit
		}

		if inp.Configs.BaseURL == "" {
			fmt.Println("The base URL shouldn't be empty. \nData will be sent to os.StdOut instead.")
		} else {
			gen.Configs.BaseURL = inp.Configs.BaseURL
		}

		if inp.Configs.BufferPostURL == "" {
			fmt.Println("The buffer post URL shouldn't be empty. \nData will be sent to os.StdOut instead.")
		} else {
			gen.Configs.BufferPostURL = inp.Configs.BufferPostURL
		}

		if inp.Configs.StreamPostURL == "" {
			fmt.Println("The stream post URL shouldn't be empty. \nData will be sent to os.StdOut instead.")
		} else {
			gen.Configs.StreamPostURL = inp.Configs.StreamPostURL
		}

	} else {
		gen.SensorCount = int32(getRandomInRange(5, 15))
		gen.LatestID1.Store(0)
		gen.Configs.DataChanSize = Configs.DataChanSize
		gen.Configs.BufferLimit = Configs.BufferLimit

		fmt.Println("The base URL shouldn't be empty. \nData will be sent to os.StdOut instead.")
		fmt.Println("The buffer post URL shouldn't be empty. \nData will be sent to os.StdOut instead.")
		fmt.Println("The stream post URL shouldn't be empty. \nData will be sent to os.StdOut instead.")
	}

	return gen, nil
}

func (g *Generator) Start() {
	beforeRoutinesCount := runtime.NumGoroutine()

	for i := 0; i < int(g.SensorCount); i++ {
		currID1 := g.LatestID1.Add(1)
		g.wg.Add(1)
		go func(id1 int32) {
			g.SpawnSensor(id1)
		}(currID1)
	}

	// We only use 1 flusher for now
	g.wg.Add(1)
	go func() {
		g.SpawnFlush()
	}()

	fmt.Print("\033[2J") // clear screen
	fmt.Print("\033[H")  // move to top

	afterRoutinesCount := runtime.NumGoroutine()
	fmt.Printf("\033[3;0H\033[2KRoutines count went from %d to %d : ", beforeRoutinesCount, afterRoutinesCount)
	if beforeRoutinesCount+int(g.SensorCount)+1 == afterRoutinesCount {
		fmt.Printf("All sensors have been spawned successfully.\n")
	} else {
		fmt.Printf("Failed to spawn all sensors.\n")
	}
}

func (g *Generator) Stop(err error) {
	g.ctxCancel(err)
	g.wg.Wait()
}

func (g *Generator) NewError(err error) {
	fmt.Printf("\033[%d;0H\033[2KErrors count : %d", g.SensorCount+10, len(g.errChan))
	fmt.Printf("\033[%d;0H\033[2KError : %s", g.SensorCount+11, err.Error())

	if len(g.errChan) > 9 {
		g.Stop(<-g.errChan)
		return
	}

	g.errChan <- err
}

func (g *Generator) GeneralMsg(msg string) {
	fmt.Printf("\033[%d;0H\033[2K%s", g.SensorCount+8, msg)
}
