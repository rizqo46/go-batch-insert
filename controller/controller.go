package controller

import (
	"go-batch-insert/database"
	"go-batch-insert/model"
	"go-batch-insert/transport"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var DB *gorm.DB = database.InitDb()

func AddSingle(c echo.Context) error {
	req := new(transport.AddSingleRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	student := &model.Student{
		Name:        req.Name,
		DateOfBirth: dateOfBirth,
		Grade:       req.Grade,
	}

	if err := DB.Create(&student).Error; err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, transport.Response{
		Succces: true,
		Message: "Success created a student",
	})
}
