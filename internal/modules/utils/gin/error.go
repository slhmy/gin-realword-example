package gin_utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidParam = ServiceError{Err: errors.New("invalid parameter"), HttpStatus: http.StatusBadRequest}
	ErrUnauthorized = ServiceError{Err: errors.New("unauthorized"), HttpStatus: http.StatusUnauthorized}

	errhttpStatusMapping = map[string]int{}
)

// Register the mapping of an error to an HTTP status code.
// If the error is a wrapped error, the base error will be used for the mapping.
func RegisterErrHttpStatusMapping(err error, httpStatus int) {
	if err == nil {
		return
	}
	fmt.Printf("RegisterErrHttpStatusMapping: %v, %d\n", err, httpStatus)
	errhttpStatusMapping[err.Error()] = httpStatus
}

type ServiceError struct {
	Err        error
	HttpStatus int
}

func (e ServiceError) Error() string {
	return e.Err.Error()
}

func (e ServiceError) Msg() string {
	return e.Err.Error()
}

// If the base error is registered, status code will be automatically set.
func WrapServiceError(err error) *ServiceError {
	if err == nil {
		return nil
	}
	if serviceErr, ok := err.(ServiceError); ok {
		return &serviceErr
	}
	httpStatus, ok := errhttpStatusMapping[err.Error()]
	if !ok {
		httpStatus = 500
	}

	return &ServiceError{
		Err:        err,
		HttpStatus: httpStatus,
	}
}

func AbortWithError(ginCtx *gin.Context, err error) {
	if err != nil {
		_ = ginCtx.Error(err)
	}
	ginCtx.Abort()
}
