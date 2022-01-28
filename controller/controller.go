package controller

import (
	"go-batch-insert/database"
	"go-batch-insert/model"
	"go-batch-insert/transport"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

var DB *gorm.DB = database.InitDb()

func AddSingle(c echo.Context) error {
	req := new(transport.AddSingleRequest)

	if err := c.Bind(req); err != nil {
		log.Errorf("error : %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dateOfBirth, err := time.Parse("01/02/2006", req.DateOfBirth)
	if err != nil {
		log.Errorf("error : %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	student := &model.Student{
		Name:        req.Name,
		DateOfBirth: dateOfBirth,
		Grade:       req.Grade,
	}

	if err := DB.Create(&student).Error; err != nil {
		log.Errorf("error : %v", err.Error())
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, transport.Response{
		Succces: true,
		Message: "Success created a student data",
	})
}

func AddBatch(c echo.Context) error {

	CsvFile, err := c.FormFile("csv_file")
	if err != nil {
		log.Errorf("error : %v", err.Error())
		return echo.ErrBadRequest
	}

	Csv, err := CsvFile.Open()
	if err != nil {
		log.Errorf("error : %v", err.Error())
		return echo.ErrBadRequest
	}
	r, err := excelize.OpenReader(Csv)
	if err != nil {
		log.Errorf("error : %v", err.Error())
		return echo.ErrBadRequest
	}

	listSheet := r.GetSheetList()

	records, err := r.GetRows(listSheet[0])
	if err != nil {
		log.Errorf("error : %v", err.Error())
		return echo.ErrBadRequest
	}

	students := make([]*model.Student, len(records)-1)

	for i, val := range records[1:] {
		dob, err := time.Parse("01/02/2006", val[1])
		if err != nil {
			log.Errorf("error : %d %v", i, err.Error())
			return echo.ErrBadRequest
		}

		grade, err := strconv.Atoi(val[2])
		if err != nil {
			log.Errorf("error : %v", err.Error())
			return echo.ErrBadRequest
		}

		students[i] = &model.Student{
			Name:        val[0],
			DateOfBirth: dob,
			Grade:       grade,
		}
	}

	numBatch := len(students)/5 + 1
	rigthIndex, lenBatch := 0, 5

	if err := DB.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < numBatch; i++ {
			switch i {
			case numBatch - 1:
				rigthIndex = int(math.Min(float64((i+1)*5), float64(len(students))))
				lenBatch = int(math.Min(float64(5), float64(rigthIndex-(i*5))))
			default:
				rigthIndex = (i + 1) * 5
			}

			if err := tx.CreateInBatches(students[i*5:rigthIndex], lenBatch).Error; err != nil {
				log.Errorf("error insert to db: %v", err.Error())
				return err
			}
		}
		return nil
	}); err != nil {
		log.Errorf("error insert to db: %v", err.Error())
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, transport.Response{
		Succces: true,
		Message: "Success to add data",
	})

}

func AddBatchUsingGoRoutine(c echo.Context) error {

	CsvFile, err := c.FormFile("csv_file")
	if err != nil {
		log.Errorf("error : %v", err.Error())
		return echo.ErrBadRequest
	}

	Csv, err := CsvFile.Open()
	if err != nil {
		log.Errorf("error : %v", err.Error())
		return echo.ErrBadRequest
	}
	r, err := excelize.OpenReader(Csv)
	if err != nil {
		log.Errorf("error : %v", err.Error())
		return echo.ErrBadRequest
	}

	go func() {
		listSheet := r.GetSheetList()

		records, err := r.GetRows(listSheet[0])
		if err != nil {
			log.Errorf("error : %v", err.Error())
		}

		students := make([]*model.Student, len(records)-1)

		for i, val := range records[1:] {
			dob, err := time.Parse("01/02/2006", val[1])
			if err != nil {
				log.Errorf("error : %d %v", i, err.Error())
			}

			grade, err := strconv.Atoi(val[2])
			if err != nil {
				log.Errorf("error : %v", err.Error())
			}

			students[i] = &model.Student{
				Name:        val[0],
				DateOfBirth: dob,
				Grade:       grade,
			}
		}

		numBatch := len(students)/5 + 1
		rigthIndex, lenBatch := 0, 5

		if err := DB.Transaction(func(tx *gorm.DB) error {
			for i := 0; i < numBatch; i++ {
				switch i {
				case numBatch - 1:
					rigthIndex = int(math.Min(float64((i+1)*5), float64(len(students))))
					lenBatch = int(math.Min(float64(5), float64(rigthIndex-(i*5))))
				default:
					rigthIndex = (i + 1) * 5
				}

				if err := tx.CreateInBatches(students[i*5:rigthIndex], lenBatch).Error; err != nil {
					log.Errorf("error insert to db: %v", err.Error())
					return err
				}
			}
			return nil
		}); err != nil {
			log.Errorf("error insert to db: %v", err.Error())
		} else {
			log.Infof("succes add %d data", len(students))
		}
	}()

	return c.JSON(http.StatusCreated, transport.Response{
		Succces: true,
		Message: "Success to add data",
	})

}
