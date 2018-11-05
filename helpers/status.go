package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
)

type Output struct {
	ReceiveFrom              int     `json:"ReceiveFrom,omitempty"`
	StatusMessage            string  `json:"StatusMessage,omitempty"`
	Bay                      int     `json:"Bay,omitempty"`
	Mode                     string  `json:"Mode,omitempty"`
	Type                     string  `json:"Type,omitempty"`
	Status                   int     `json:"Status,omitempty"`
	Name                     string  `json:"Name,omitempty"`
	DPS                      int     `json:"DPS,omitempty"`
	HPD                      int     `json:"HPD,omitempty"`
	HDCP                     int     `json:"HDCP,omitempty"`
	HasSignal                bool    `json:"HasSignal,omitempty"`
	LinkStatus               string  `json:"LinkStatus,omitempty"`
	FirmwareVersion          string  `json:"FirmwareVersion,omitempty"`
	FirmwareVersionAvailable bool    `json:"FirmwareVersionAvailable"`
	SupportSB                bool    `json:"SupportSB,omitempty"`
	Result                   bool    `json:"Result,omitempty"`
	ErrorMessage             *string `json:"ErrorMessage,omitempty"`
}

type Input struct {
	TransmissionNodes []int  `json:"TransmissionNodes"`
	StatusMessage     string `json:"StatusMessage"`
	Bay               int    `json:"Bay"`
	Mode              string `json:"Mode"`
	Type              string `json:"Type"`
	Status            int    `json:"Status"`
	Name              string `json:"Name"`
	EdidProfile       int    `json:"EdidProfile"`
	DPS               int    `json:"DPS"`
	CECVersion        int    `json:"CEC_version"`
	HPD               int    `json:"HPD"`
	HDCP              int    `json:"HDCP"`
	HasSignal         bool   `json:"HasSignal"`
	Result            bool   `json:"Result"`
}

type Port struct {
	Bay        int    `json:"Bay,omitempty"`
	Mode       string `json:"Mode,omitempty"`
	Type       string `json:"Type,omitempty"`
	Status     int    `json:"Status,omitempty"`
	Name       string `json:"Name,omitempty"`
	DPS        int    `json:"DPS,omitempty"`
	CECVersion int    `json:"CEC_version,omitempty"`
}

type PortList struct {
	Result bool   `json:"Result,omitempty"`
	Ports  []Port `json:"Ports,omitempty"`
}

func GetCurrentInputs(address string) (map[string]string, error) {

	outputMap := make(map[string]string)

	response, err := http.Get(fmt.Sprintf("http://%s/Port/List", address))
	if err != nil {
		return make(map[string]string), err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return make(map[string]string), err
	}

	var ports PortList
	err = json.Unmarshal(body, &ports)
	if err != nil {
		return make(map[string]string), err
	}

	for _, port := range ports.Ports {
		if port.Mode == "Output" {

			in, err := GetInputByOutputPort(address, port.Bay)
			if err != nil {
				return make(map[string]string), err
			}
			input := in.Input

			outputMap[fmt.Sprintf("%v", port.Bay)] = input
		}
	}

	return outputMap, nil
}

//returns a string representing the name of the source displayed on the output port
//@param bay the number of the physical bay on the device
func getInputInfoByOutputPort(address string, bay int) (Output, error) {

	log.L.Infof("Querying input port for output bay %v", bay+1)
	var output Output

	//make a call to get the input source
	response, err := http.Get(fmt.Sprintf("http://%s/Port/Details/Output/%v", address, bay))
	if err != nil {
		return output, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return output, err
	}

	err = json.Unmarshal(body, &output)
	if err != nil {
		return output, err
	}
	if output.ErrorMessage != nil {
		return output, errors.New(*output.ErrorMessage)
	}
	return output, nil
}

//returns a string representing the name of the source displayed on the output port
//@param bay the number of the physical bay on the device
func GetInputNameByOutputPort(address string, bay int) (string, error) {
	output, err := getInputInfoByOutputPort(address, bay)
	if err != nil {
		return "", err
	}

	var input Input
	log.L.Infof("Querying input bay %v", output.ReceiveFrom+1)

	//make a call based on the bay the output recieves from
	response, err := http.Get(fmt.Sprintf("http://%s/Port/Details/Input/%v", address, output.ReceiveFrom))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		return "", err
	}

	return input.Name, nil
}

func GetInputByOutputPort(address string, bay int) (status.Input, error) {
	output, err := getInputInfoByOutputPort(address, bay)
	if err != nil {
		return status.Input{}, err
	}

	return status.Input{Input: fmt.Sprintf("%v", output.ReceiveFrom)}, nil
}
