package transport

type Response struct {
	Succces bool   `json:"succes"`
	Message string `json:"message"`
}

type AddSingleRequest struct {
	Name        string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
	Grade       int    `json:"grade"`
}
