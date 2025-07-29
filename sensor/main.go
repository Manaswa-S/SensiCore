package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sensor/internal"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {

	flowChan := make(chan os.Signal, 1)
	signal.Notify(flowChan, syscall.SIGINT, syscall.SIGTERM)

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	baseURL, exists := os.LookupEnv("BACKEND_BASE_URL")
	if !exists {
		fmt.Println("Backend base url does not exist in .env.")
		return
	}
	buffpostPath, exists := os.LookupEnv("FLUSH_DATA_PATH_BUFFERED")
	if !exists {
		fmt.Println("Flush data path buffered does not exist in .env.")
		return
	}
	strmpostPath, exists := os.LookupEnv("FLUSH_DATA_PATH_STREAMED")
	if !exists {
		fmt.Println("Flush data path streamed does not exist in .env.")
		return
	}

	stream := false
	sensors := 3

	if streamStr, exists := os.LookupEnv("STREAM_DATA"); exists {
		if streamStr == "true" {
			stream = true
		}
	}
	if sensorsStr, exists := os.LookupEnv("SENSORS_COUNT"); exists {
		sensors, err = strconv.Atoi(sensorsStr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if f := flag.Lookup("stream"); f != nil {
		stream = *flag.Bool("stream", false, "To stream data or not")
	}
	if f := flag.Lookup("sensors"); f != nil {
		sensors = *flag.Int("sensors", 6, "Number of sensors to spawn")
	}

	flag.Parse()

	generator, err := internal.NewSensorGenerator(&internal.Generator{
		StreamFlush: stream,
		SensorCount: int32(sensors),
		Configs: internal.Config{
			// DataChanSize: 100,
			// BufferLimit:  70,
			BaseURL:       baseURL,
			BufferPostURL: baseURL + buffpostPath,
			StreamPostURL: baseURL + strmpostPath,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// generator.LatestID1.Store(122) // way to override LatestID1

	generator.Start()

	<-flowChan

	generator.Stop(nil)
}
