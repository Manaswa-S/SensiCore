package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"sensicore/internal/dto"
	sqlc "sensicore/internal/sqlc/generate"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Service struct {
	SQLDB   *sql.DB
	Queries *sqlc.Queries
}

func NewService(sqldb *sql.DB, queries *sqlc.Queries) *Service {
	return &Service{
		SQLDB:   sqldb,
		Queries: queries,
	}
}

func (s *Service) PostData(eCtx echo.Context, buffer *[]dto.SensorDataReq) (*dto.SensorDataResp, error) {
	ctx := eCtx.Request().Context()

	sqlTx, err := s.SQLDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	txQueries := s.Queries.WithTx(sqlTx)

	// TODO: maybe we need a defer rollback here ??

	successfulCnt := int64(len(*buffer))

	for _, d := range *buffer {
		err = txQueries.InsertSensorData(ctx, sqlc.InsertSensorDataParams{
			Value:  d.Value,
			Unit:   sql.NullString{String: d.Unit, Valid: true},
			Id1:    d.ID1,
			Id2:    d.ID2,
			ReadAt: d.ReadAt,
		})
		if err != nil {
			fmt.Println(err) // TODO: silent error drop
			successfulCnt--
		}
	}
	if err := sqlTx.Commit(); err != nil {
		return nil, err
	}

	return &dto.SensorDataResp{
		SuccessfulCnt: successfulCnt,
	}, nil
}

func (s *Service) PostDataStream(eCtx echo.Context) (*dto.SensorDataResp, error) {
	ctx := eCtx.Request().Context()
	// defer eCtx.Request().Body.Close()

	decoder := json.NewDecoder(eCtx.Request().Body)
	successfulCnt := int64(0)

	for {

		data := new(dto.SensorDataReq)

		if err := decoder.Decode(data); err != nil {
			if err == io.EOF {
				return &dto.SensorDataResp{
					SuccessfulCnt: successfulCnt,
				}, nil
			}
			return nil, err
		}
		err := s.Queries.InsertSensorData(ctx, sqlc.InsertSensorDataParams{
			Value:  data.Value,
			Unit:   sql.NullString{String: data.Unit, Valid: true},
			Id1:    data.ID1,
			Id2:    data.ID2,
			ReadAt: data.ReadAt,
		})
		if err != nil {
			return nil, err
		}
		successfulCnt++
	}
}

func (s *Service) GetData(eCtx echo.Context, id1Str, id2Str, startStr, endStr, limitStr, offsetStr string) ([]sqlc.GetAllSensorDataRow, error) {
	ctx := eCtx.Request().Context()

	var iD1 int64
	var iD2 string
	var start int64
	var end int64

	limit := int64(25)
	offset := int64(0)

	var err error
	iD2 = id2Str

	if id1Str != "" {
		iD1, err = strconv.ParseInt(id1Str, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	iD2 = id2Str

	if startStr != "" {
		start, err = strconv.ParseInt(startStr, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if endStr != "" {
		end, err = strconv.ParseInt(endStr, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if limitStr != "" {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return nil, err
		}
		if limit <= 0 {
			limit = 25
		}
	}

	if offsetStr != "" {
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			return nil, err
		}
		if offset < 0 {
			offset = 0
		}
	}

	startTime := time.Unix(start, 0).UTC()
	endTime := time.Unix(end, 0).UTC()

	data, err := s.Queries.GetAllSensorData(ctx, sqlc.GetAllSensorDataParams{
		Id1:      int32(iD1),
		Column2:  iD1,
		Id2:      iD2,
		Column4:  iD2,
		ReadAt:   startTime,
		Column6:  startTime,
		Column7:  time.Unix(0, 0).UTC(), // Zero time is passed instead of string comparison
		ReadAt_2: endTime,
		Column9:  endTime,
		Column10: time.Unix(0, 0).UTC(),
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
