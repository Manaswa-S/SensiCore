package handlers

import (
	"fmt"
	"net/http"
	"sensicore/internal/dto"

	"github.com/labstack/echo/v4"
)

func (h *Handler) PostData(eCtx echo.Context) error {

	buffer := new([]dto.SensorDataReq)

	err := eCtx.Bind(buffer)
	if err != nil {
		fmt.Println(err)
		return eCtx.NoContent(http.StatusBadRequest)
	}

	resp, err := h.Service.PostData(eCtx, buffer)
	if err != nil {
		fmt.Println(err)
		return eCtx.NoContent(http.StatusBadRequest)
	}

	return eCtx.JSON(http.StatusCreated, resp)
}

func (h *Handler) PostDataStream(eCtx echo.Context) error {

	resp, err := h.Service.PostDataStream(eCtx)
	if err != nil {
		fmt.Println(err)
		return eCtx.NoContent(http.StatusBadRequest)
	}

	return eCtx.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetData(eCtx echo.Context) error {

	id1 := eCtx.QueryParam("id1")
	id2 := eCtx.QueryParam("id2")
	start := eCtx.QueryParam("start")
	end := eCtx.QueryParam("end")

	limit := eCtx.QueryParam("limit")
	offset := eCtx.QueryParam("offset")

	resp, err := h.Service.GetData(eCtx, id1, id2, start, end, limit, offset)
	if err != nil {
		fmt.Println(err)
		return eCtx.NoContent(http.StatusBadRequest)
	}

	return eCtx.JSON(http.StatusOK, resp)
}
