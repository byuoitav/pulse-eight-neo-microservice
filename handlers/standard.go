package handlers

import (
	"net/http"

	"github.com/byuoitav/pulse-eight-neo-microservice/helpers"
	"github.com/labstack/echo"
)

func PowerOn(context echo.Context) error {
	return echo.JSON(http.StatusOK, nil)
}

func Standby(context echo.Context) error {
	return echo.JSON(http.StatusOK, nil)
}

func SwitchInput(context echo.Context) error {

	input := context.Param("input")
	output := context.Param("output")
	address := context.Param("address")

	err := helpers.SwitchInput(address, input, output)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context > JSON(http.StatusOK, "Success")
}
