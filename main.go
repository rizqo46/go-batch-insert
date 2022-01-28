package main

import (
	"go-batch-insert/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/add/single", controller.AddSingle)
	e.POST("/add/batch", controller.AddBatch)
	e.POST("/add/batch/async", controller.AddBatchUsingGoRoutine)

	e.Logger.Fatal(e.Start(":1323"))
}
