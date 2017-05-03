package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Output struct {
	ReceiveFrom              int    `json:"ReceiveFrom",omitempty`
	StatusMessage            string `json:"StatusMessage",omitempty`
	Bay                      int    `json:"Bay"omitempty`
	Mode                     string `json:"Mode",omitempty`
	Type                     string `json:"Type",omitempty`
	Status                   int    `json:"Status"omitempty`
	Name                     string `json:"Name",omitempty`
	DPS                      int    `json:"DPS",omitempty`
	HPD                      int    `json:"HPD",omitempty`
	HDCP                     int    `json:"HDCP",omitempty`
	HasSignal                bool   `json:"HasSignal",omitempty`
	LinkStatus               string `json:"LinkStatus",omitempty`
	FirmwareVersion          string `json:"FirmwareVersion",omitempty`
	FirmwareVersionAvailable bool   `json:"FirmwareVersionAvailable"`
	SupportSB                bool   `json:"SupportSB",omitempty`
	Result                   bool   `json:"Result",omitempty`
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
	CEC_version       int    `json:"CEC_version"`
	HPD               int    `json:"HPD"`
	HDCP              int    `json:"HDCP"`
	HasSignal         bool   `json:"HasSignal"`
	Result            bool   `json:"Result"`
}

type Port struct {
	Bay         int    `json:"Bay",omitempty`
	Mode        string `json:"Mode",omitempty`
	Type        string `json:"Type",omitempty`
	Status      int    `json:"Status",omitempty`
	Name        string `json:"Name",omitempty`
	DPS         int    `json:"DPS",omitempty`
	CEC_version int    `json:"CEC_version",omitempty`
}

type PortList struct {
	Result bool   `json:"Result",omitempty`
	Ports  []Port `json:"Ports",omitempty`
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

			input, err := GetInputByOutputPort(address, port.Bay)
			if err != nil {
				return make(map[string]string), err
			}

			outputMap[port.Name] = input
		}
	}

	return outputMap, nil
}

//returns a string representing the name of the source displayed on the output port
//@param bay the number of the physical bay on the device
func GetInputByOutputPort(address string, bay int) (string, error) {

	//make a call to get the input source
	response, err := http.Get(fmt.Sprintf("http://%s/Port/Details/Output/%v", address, bay))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var output Output
	err = json.Unmarshal(body, &output)
	if err != nil {
		return "", err
	}

	//make a call based on the bay the output recieves from
	response, err = http.Get(fmt.Sprintf("http://%s/Port/Details/Input/%v", address, output.ReceiveFrom))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var input Input
	err = json.Unmarshal(body, &input)
	if err != nil {
		return "", err
	}

	return input.Name, nil
}
