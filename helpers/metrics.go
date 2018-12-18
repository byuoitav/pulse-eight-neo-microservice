package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
)

// NeoHardwareResponse is the struct that the Neo switcher reports information in
type NeoHardwareResponse struct {
	BoardRev        int    `json:"BoardRev"`
	MACAddress      string `json:"MAC"`
	ModelName       string `json:"Model"`
	Result          bool   `json:"Result"`
	SerialNumber    string `json:"Serial"`
	StatusCode      int    `json:"Status"`
	StatusMessage   string `json:"StatusMessage"`
	VID             string `json:"VID"`
	FirmwareVersion string `json:"Version"`
}

// GetHardwareInfo performs the functions necessary to build the hardware information struct and returns it.
func GetHardwareInfo(address string) (structs.HardwareInfo, *nerr.E) {
	var toReturn structs.HardwareInfo

	resp, err := http.Get(fmt.Sprintf("http://%s/System/Details", address))
	if err != nil {
		msg := fmt.Sprintf("failed to get a response from %s", address)
		log.L.Errorf("%s : %s", msg, err.Error())
		return toReturn, nerr.Translate(err).Add(msg)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("failed to read the response from %s", address)
		log.L.Errorf("%s : %s", msg, err.Error())
		return toReturn, nerr.Translate(err).Add(msg)
	}

	var neo NeoHardwareResponse
	err = json.Unmarshal(body, &neo)
	if err != nil {
		msg := fmt.Sprintf("failed to unmarshal the response from %s", address)
		log.L.Errorf("%s : %s", msg, err.Error())
		return toReturn, nerr.Translate(err).Add(msg)
	}

	// get the hostname
	addr, e := net.LookupAddr(address)
	if e != nil {
		toReturn.Hostname = address

		ip, _ := net.LookupHost(address)
		if len(ip) > 0 {
			toReturn.NetworkInfo.IPAddress = ip[0]
		}
	} else {
		toReturn.Hostname = strings.Trim(addr[0], ".")
		toReturn.NetworkInfo.IPAddress = address
	}

	toReturn.ModelName = neo.ModelName
	toReturn.SerialNumber = neo.SerialNumber
	toReturn.FirmwareVersion = neo.FirmwareVersion
	toReturn.NetworkInfo.MACAddress = neo.MACAddress

	if !strings.EqualFold(neo.StatusMessage, "Healthy") || neo.StatusCode != 0 {
		toReturn.ErrorStatus = append(toReturn.ErrorStatus, neo.StatusMessage)
	}

	return toReturn, nil
}

// GetActiveSignal takes an address and a port and determines if there is an active signal on that port.
func GetActiveSignal(address, port string) (structs.ActiveSignal, *nerr.E) {
	var toReturn structs.ActiveSignal

	return toReturn, nil
}
