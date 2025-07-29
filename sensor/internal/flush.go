package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (g *Generator) SpawnFlush() {
	defer g.wg.Done()

	if g.StreamFlush {
		g.streamFlush()
	} else {
		g.bufferFlush()
	}
}

func (g *Generator) streamFlush() {

	encoder := g.acquireStreamConn()

	totalFlushCnt := int64(0)
	successfulFlushCnt := int64(0)

	for {

		fmt.Printf("\033[5;0H\033[2KFlushing : Flushed : %5d : Successful : %8d",
			totalFlushCnt, successfulFlushCnt)

		select {
		case <-g.ctx.Done():
			// clean up, send all data first
			dChnPtr := g.dataChan.Load()
			for len(*dChnPtr) > 0 {
				if err := encoder.Encode(<-*dChnPtr); err != nil {
					if err == io.ErrClosedPipe {
						encoder = g.acquireStreamConn()
					} else {
						g.NewError(err)
					}
				}
			}
			return
		case data := <-*g.dataChan.Load():
			if err := encoder.Encode(data); err != nil {
				if err == io.ErrClosedPipe {
					encoder = g.acquireStreamConn()
				} else {
					g.NewError(err)
				}
			} else {
				totalFlushCnt++
			}
		}
	}
}

func (g *Generator) acquireStreamConn() *json.Encoder {

	pipeReader, pipeWriter := io.Pipe()

	req, err := http.NewRequest("POST", g.Configs.StreamPostURL, pipeReader)
	if err != nil {
		g.NewError(err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")

	go func(pipeWriter *io.PipeWriter) {
		go func() {
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				g.NewError(err)
				return
			}
			g.GeneralMsg("Acquired stream connection.")
			defer resp.Body.Close()
		}()

		time.Sleep(15 * time.Second) // The connection is refreshed every 15 seconds.

		pipeWriter.Close()
	}(pipeWriter)

	return json.NewEncoder(pipeWriter)
}

func (g *Generator) bufferFlush() {

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	totalFlushCnt := int64(0)
	successfulFlushCnt := int64(0)

	for {

		fmt.Printf("\033[5;0H\033[2KFlushing : Flushed : %5d : Successful : %8d",
			totalFlushCnt, successfulFlushCnt)

		select {
		case <-g.ctx.Done():
			// clean up, send all data first
			g.flushBuffer()
			return
		case <-ticker.C:
			dChnPtr := g.dataChan.Load()
			if len(*dChnPtr) >= int(Configs.BufferLimit) {

				// Swap the channels
				temp := g.dataChan.Swap(g.altDataChan.Load())
				g.altDataChan.Store(temp)
				g.GeneralMsg("Data channels swapped.")

				tfC, sfC := g.flushBuffer()
				totalFlushCnt += tfC
				successfulFlushCnt += sfC
			}
		}
	}
}

func (g *Generator) flushBuffer() (totalFlushCnt, successfulFlushCnt int64) {
	adChnPtr := g.altDataChan.Load()

	buffer := make([]*SensorData, 0)
	totalFlushCnt = 0

	for len(*adChnPtr) > 0 {
		buffer = append(buffer, <-*adChnPtr)
		totalFlushCnt++
	}

	body, err := json.Marshal(buffer)
	if err != nil {
		g.NewError(err)
		return
	}

	req, err := http.NewRequest("POST", g.Configs.BufferPostURL, bytes.NewBuffer(body))
	if err != nil {
		g.NewError(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		g.NewError(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			g.NewError(err)
			return
		}
		var resBody SensorDataResp
		if err := json.Unmarshal(bodyBytes, &resBody); err != nil {
			g.NewError(err)
			return
		}

		successfulFlushCnt = resBody.SuccessfulCnt
	}

	return
}
