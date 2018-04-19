package resource

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// BadRequest http 400 error
func BadRequest(payload interface{}) error {
	return echo.NewHTTPError(http.StatusBadRequest, processPayload(payload))
}

// NotFound http 404 error
func NotFound(payload interface{}) error {
	return echo.NewHTTPError(http.StatusNotFound, processPayload(payload))
}

// InternalServerError http 500 error
func InternalServerError(payload interface{}) error {
	return echo.NewHTTPError(http.StatusInternalServerError, processPayload(payload))
}

// Unauthorized http 401 status
func Unauthorized() error {
	return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
}

// Created http 201 response
func Created(c echo.Context, id uint) error {
	c.Response().Header().Set("Location",
		fmt.Sprintf("%s/%d", c.Request().URL.String(), id))
	return c.NoContent(http.StatusCreated)
}

func processPayload(payload interface{}) interface{} {
	if err, isError := payload.(error); isError {
		return err.Error
	}
	return payload
}
