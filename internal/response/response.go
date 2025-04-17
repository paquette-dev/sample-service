package response

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error string `json:"error"`
}

// JSONSuccessResponse returns a success response
func JSONSuccessResponse(ctx echo.Context, message string, data interface{}) error {
	return ctx.JSON(http.StatusOK, SuccessResponse{
		Message: message,
		Data: data,
	})
}

// JSONErrorResponse returns an error response
func JSONErrorResponse(ctx echo.Context, message string, error string) error {
	return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
		Message: message,
		Error: error,
	})
}
