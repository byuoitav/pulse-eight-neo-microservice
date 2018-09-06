package handlers

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/byuoitav/pulse-eight-neo-microservice/helpers"
	"github.com/labstack/echo"
)

const lockTime = 128 * time.Millisecond

func SwitchInput(context echo.Context) error {
	log.SetLevel("debug")
	input := context.Param("input")
	output := context.Param("output")
	address := context.Param("address")

	log.L.Infof("Routing %v to %v on %v", input, output, address)

	lock(address)
	defer unlock(address)

	err := helpers.SwitchInput(address, input, output)
	if err != nil {
		log.L.Warnf("unable to switch inputs: %v", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	log.L.Infof("Success")
	returnVal := input + ":" + output

	return context.JSON(http.StatusOK, status.Input{Input: returnVal})
}

func GetCurrentInput(context echo.Context) error {
	address := context.Param("address")

	lock(address)
	defer unlock(address)

	inputs, err := helpers.GetCurrentInputs(address)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, inputs)
}

func GetInputByPort(context echo.Context) error {
	address := context.Param("address")
	port := context.Param("port")
	bay, err := strconv.Atoi(port)
	if err != nil || bay < 0 {
		return context.JSON(http.StatusBadRequest, "Error! Port parameter must be zero or greater")
	}

	lock(address)
	defer unlock(address)

	input, err := helpers.GetInputByOutputPort(address, bay)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	input.Input = input.Input + ":" + port

	return context.JSON(http.StatusOK, input)
}

var (
	delayMap map[string]*sync.Mutex
	mapInit  sync.Once
	mapLock  sync.RWMutex
)

func lock(address string) {
	mapInit.Do(func() {
		delayMap = make(map[string]*sync.Mutex)
	})

	mapLock.RLock()
	lock := delayMap[address]
	mapLock.RUnlock()

	if lock == nil {
		mapLock.Lock()

		lock = &sync.Mutex{}
		delayMap[address] = lock

		mapLock.Unlock()
	}

	log.L.Debugf("Waiting for lock on address %v", address)
	lock.Lock()
	log.L.Debugf("Received lock on %v", address)
}

func unlock(address string) {
	mapLock.RLock()
	lock := delayMap[address]
	mapLock.RUnlock()

	if lock == nil {
		return
	}

	go func() {
		log.L.Debugf("Unlocking %v in %v", address, lockTime)

		time.Sleep(lockTime)
		lock.Unlock()

		log.L.Debugf("Unlocked %v", address)
	}()
}
