package dto

import "time"

type SensorDataReq struct {
	Value  float64   `json:"value"`
	Unit   string    `json:"unit"`
	ID1    int32     `json:"id1"`
	ID2    string    `json:"id2"`
	ReadAt time.Time `json:"readat"`
}

type SensorDataResp struct {
	SuccessfulCnt int64 `json:"successfulcnt"`
}
