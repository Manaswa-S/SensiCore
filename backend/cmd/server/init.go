package server

import (
	"fmt"
	"os"
	"sensicore/cmd/db"
	"sensicore/internal/handlers"
	"sensicore/internal/services"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitHTTPServer(ds *db.DataStore) error {

	e := echo.New()

	e.Server.ReadTimeout = time.Duration(30) * time.Second

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "[GIN] ${time_custom} | ${status} | ${latency_human} | ${remote_ip} | ${method} ${uri}\n",
		CustomTimeFormat: "2006/01/02 - 15:04:05",
	}))
	e.Use(middleware.Recover())

	err := initRoutes(e, ds)
	if err != nil {
		return err
	}

	go func() {
		err := e.Start(":" + os.Getenv("PORT"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	return nil
}

func initRoutes(e *echo.Echo, ds *db.DataStore) error {

	// < Init Request logger middleware directly on router.
	// requestsLoggerFile, err := os.OpenFile(configs.Paths.RequestLoggerFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	return err
	// }
	// router.Use(middlewares.Logger(requestsLoggerFile))
	// >

	sqlDB := ds.SQLDB
	queries := ds.Queries

	publicGrp := e.Group("/public")

	service := services.NewService(sqlDB, queries)
	handler := handlers.NewHandler(service)
	handler.RegisterRoutes(publicGrp)

	return nil
}
