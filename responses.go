package bubucore

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ErrorResponse sends Error as gin response
func ErrorResponse(ctx *gin.Context, err *Error) {
	status := http.StatusInternalServerError
	if err.Code >= 400 && err.Code <= 599 {
		status = err.Code
	}
	ctx.JSON(
		status,
		err,
	)
}

// ErrorResponseE send standard error as response
func ErrorResponseE(ctx *gin.Context, err error, status int) {
	ErrorResponse(ctx, NewError(status, err.Error()))
}

// ErrorResponseS sends error string as response
func ErrorResponseS(ctx *gin.Context, msg string, status int) {
	ErrorResponse(ctx, NewError(status, msg))
}
