package main

import (
	"net/http"

	"github.com/byuoitav/common"
	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/pulse-eight-neo-microservice/handlers"
	"github.com/labstack/echo"
)

func main() {
	port := ":8011"
	router := common.NewRouter()

	// Use the `rotuer` routing group to require authentication

	router.GET("/", echo.WrapHandler(http.HandlerFunc(hateoas.RootResponse)))

	//Functionality endpoints
	router.GET("/:address/input/:input/:output", handlers.SwitchInput)

	//Status endpoints
	router.GET("/:address/input/map", handlers.GetCurrentInput)
	router.GET("/:address/input/get/:port", handlers.GetInputByPort)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
