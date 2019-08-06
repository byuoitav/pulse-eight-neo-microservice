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

//SwitchInput .
func SwitchInput(context echo.Context) error {
	input := context.Param("input")
	output := context.Param("output")
	address := context.Param("address")

	log.L.Infof("Routing %v to %v on %v", input, output, address)

	unlock := lock(address)
	defer unlock()

	err := helpers.SwitchInput(address, input, output)
	if err != nil {
		log.L.Warnf("unable to switch inputs: %v", err)
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	log.L.Infof("Success")
	returnVal := input + ":" + output

	return context.JSON(http.StatusOK, status.Input{Input: returnVal})
}

// GetCurrentInput returns what the current input being shown is
func GetCurrentInput(context echo.Context) error {
	address := context.Param("address")

	unlock := lock(address)
	defer unlock()

	inputs, err := helpers.GetCurrentInputs(address)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, inputs)
}

// GetInputByPort will return the input for a specific port
func GetInputByPort(context echo.Context) error {
	address := context.Param("address")
	port := context.Param("port")

	bay, err := strconv.Atoi(port)
	if err != nil || bay < 0 {
		return context.JSON(http.StatusBadRequest, "Error! Port parameter must be zero or greater")
	}

	unlock := lock(address)
	defer unlock()

	input, err := helpers.GetInputByOutputPort(address, bay)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	input.Input = input.Input + ":" + port

	return context.JSON(http.StatusOK, input)
}

// IssueReboot will reboot a pulse 8 switcher
func IssueReboot(context echo.Context) error {
	address := context.Param("address")

	// Send the reboot command
	err := helpers.RebootNeo(address)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, status.Power{Power: "Rebooting"})
}

var (
	delayMap map[string]*sync.Mutex
	mapInit  sync.Once
	mapLock  sync.RWMutex
)

func lock(address string) func() {
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

	return func() {
		go func() {
			log.L.Debugf("Unlocking %v in %v", address, lockTime)

			time.Sleep(lockTime)
			lock.Unlock()

			log.L.Debugf("Unlocked %v", address)
		}()
	}
}
