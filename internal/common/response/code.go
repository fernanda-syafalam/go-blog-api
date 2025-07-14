package response

import "net/http"

type Code string

const (
	Success Code = "20000"
	Created Code = "20100"

	BadRequest          Code = "40000"
	InvalidRequest      Code = "40001"
	InvalidID           Code = "40003"
	ResourceNotFound    Code = "40004"
	ResouceAlreadyExist Code = "40005"
	CheckinInvalid      Code = "40002"

	Unauthorized  Code = "40100"
	InvalidToken  Code = "40101"
	TokenNotFound Code = "40102"

	Forbidden Code = "40300"

	ServerError        Code = "50000"
	Timeout            Code = "50400"
	ServiceUnavailable Code = "50300"
)

var (
	codeMap = map[Code]string{
		Success: "Success",
		Created: "Created",

		BadRequest:       "Bad or invalid request",
		InvalidRequest:   "Invalid Request",
		InvalidID:        "Invalid ID",
		ResourceNotFound: "Resource Not Found",
		ResouceAlreadyExist: "Resource Already Exist",
		CheckinInvalid:   "Ticket Already Used or Expired",

		Unauthorized:  "Unauthorized Token",
		InvalidToken:  "Invalid Token",
		TokenNotFound: "Token Not Found",

		Forbidden: "Forbidden",

		ServerError:        "Internal Server Error",
		Timeout:            "Gateway Timeout",
		ServiceUnavailable: "Service Unavailable",
	}

	codeHTTPMap = map[Code]int{
		Success: http.StatusOK,
		Created: http.StatusCreated,

		BadRequest:       http.StatusBadRequest,
		InvalidRequest:   http.StatusUnprocessableEntity,
		InvalidID:        http.StatusBadRequest,
		ResourceNotFound: http.StatusNotFound,
		CheckinInvalid:   http.StatusBadRequest,

		Unauthorized:  http.StatusUnauthorized,
		InvalidToken:  http.StatusUnauthorized,
		TokenNotFound: http.StatusUnauthorized,

		Forbidden: http.StatusForbidden,

		Timeout:            http.StatusGatewayTimeout,
		ServerError:        http.StatusInternalServerError,
		ServiceUnavailable: http.StatusServiceUnavailable,
	}
)

func (c Code) GetMessage() string {
	return codeMap[c]
}

func (c Code) GetHTTPCode() int {
	return codeHTTPMap[c]
}

func (c Code) GetCode() string {
	return string(c)
}
