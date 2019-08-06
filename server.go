package main

import (
	"net/http"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/v2/auth"
	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/pulse-eight-neo-microservice/handlers"
	"github.com/labstack/echo"
)

func main() {
	port := ":8011"
	router := common.NewRouter()

	// Use the `rotuer` routing group to require authentication
	write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	read := router.Group("", auth.AuthorizeRequest("read-state", "room", auth.LookupResourceFromAddress))

	router.GET("/", echo.WrapHandler(http.HandlerFunc(hateoas.RootResponse)))

	//Functionality endpoints
	write.GET("/:address/input/:input/:output", handlers.SwitchInput)
	write.GET("/:address/reboot", handlers.IssueReboot)

	//Status endpoints
	read.GET("/:address/input/map", handlers.GetCurrentInput)
	read.GET("/:address/input/get/:port", handlers.GetInputByPort)
	read.GET("/:address/hardware", handlers.GetHardwareInfo)
	read.GET("/:address/active/:port", handlers.GetActiveSignalByPort)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
