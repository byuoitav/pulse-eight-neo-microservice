package main

import (
	"net/http"

	"github.com/byuoitav/authmiddleware"
	"github.com/byuoitav/common"
	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/pulse-eight-neo-microservice/handlers"
	"github.com/labstack/echo"
)

func main() {
	port := ":8011"
	router := common.NewRouter()

	// Use the `secure` routing group to require authentication
	secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	router.GET("/", echo.WrapHandler(http.HandlerFunc(hateoas.RootResponse)))

	//Functionality endpoints
	secure.GET("/:address/input/:input/:output", handlers.SwitchInput)

	//Status endpoints
	secure.GET("/:address/input/map", handlers.GetCurrentInput)
	secure.GET("/:address/input/get/:port", handlers.GetInputByPort)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
