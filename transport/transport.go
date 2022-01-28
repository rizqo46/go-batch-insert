package transport

import "mime/multipart"

type Response struct {
	Succces bool   `json:"succes"`
	Message string `json:"message"`
}

type AddSingleRequest struct {
	Name        string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
	Grade       int    `json:"grade"`
}

type AddBatchRequest struct {
	CsvFile *multipart.FileHeader `json:"csv_file" form:"csv_file"`
}
