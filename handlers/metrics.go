package handlers

import (
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/pulse-eight-neo-microservice/helpers"
	"github.com/labstack/echo"
)

// GetHardwareInfo returns the hardware information about the device
func GetHardwareInfo(context echo.Context) error {
	address := context.Param("address")

	hardware, err := helpers.GetHardwareInfo(address)
	if err != nil {
		log.L.Error(err.Error())
		return context.JSON(http.StatusInternalServerError, err.String())
	}

	return context.JSON(http.StatusOK, hardware)
}

// GetActiveSignalByPort checks to see if a port has an active signal
func GetActiveSignalByPort(context echo.Context) error {
	address := context.Param("address")
	port := context.Param("port")

	active, err := helpers.GetActiveSignal(address, port)
	if err != nil {
		log.L.Error(err.Error())
		return context.JSON(http.StatusInternalServerError, err.String())
	}

	return context.JSON(http.StatusOK, active)
}
