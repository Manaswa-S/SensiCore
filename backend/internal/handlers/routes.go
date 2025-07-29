package handlers

import (
	"sensicore/internal/services"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) RegisterRoutes(routeGrp *echo.Group) {

	routeGrp.POST("/data", h.PostData)
	routeGrp.POST("/data/stream", h.PostDataStream)

	// Useful Query Params
	// id1 ... filters on given id1
	// id2 ... filters on given id2
	// start ... filters on read_at >= start i.e. events that occured after start
	// end ... filters on read_at <= end i.e. events that occured before end
	// limit ... limits number of events returned, default is 15
	// offset ... ignores (does not return) the given number of events, from the top, default is 0
	// Data is generally sorted on the time it was generated, read_at.
	routeGrp.GET("/data", h.GetData)
}

// [ID1] , all rows with ID1
// [ID2] , all rows with ID2
// [start] , all rows with read_at >= start
// [end] , all rows with read_at <= end
// [ID1 ID2] , all rows with ID1 and ID2
// [ID1 start] , all rows with ID1, read_at >= start
// [ID1 end] , all rows with ID1, read_at <= end
// [ID2 start] , all rows with ID2, read_at >= start
// [ID2 end] , all rows with ID1, read_at <= end
// [start end] , all rows with start =< read_at <= end
// [ID1 ID2 start] , all rows with ID1 and ID2, read_at >= start
// [ID1 ID2 end] , , all rows with ID1 and ID2, read_at <= end
// [ID1 start end] , all rows with ID1, start =< read_at <= end
// [ID2 start end] , all rows with ID2, start =< read_at <= end
// [ID1 ID2 start end] , all rows with ID1 and ID2, start =< read_at <= end
