package helpers

import(
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Output struct {
	ReceiveFrom int `json:"ReceiveFrom",omitempty`
	StatusMessage string `json:"StatusMessage",omitempty`
	Bay int `json:"Bay"omitempty`
	Mode string `json:"Mode",omitempty`
	Type string `json:"Type",omitempty`
	Status int `json:"Status"omitempty`
	Name string `json:"Name",omitempty`
	DPS int `json:"DPS",omitempty`
	HPD int `json:"HPD",omitempty`
	HDCP int `json:"HDCP",omitempty`
	HasSignal bool `json:"HasSignal",omitempty`
	LinkStatus string `json:"LinkStatus",omitempty`
	FirmwareVersion string `json:"FirmwareVersion",omitempty`
	FirmwareVersionAvailable bool `json:"FirmwareVersionAvailable"`
	SupportSB bool `json:"SupportSB",omitempty`
	Result bool `json:"Result",omitempty`
}

func GetCurrentInputs(address string) (map[int]int, error) {

	routerMap := make(map[int]int)
	const NUM_OUTPUTS = 4

	for i := 0; i < NUM_OUTPUTS - 1; i++ {
		response, err := http.Get(fmt.Sprintf("http://%s/Port/Details/Output/%v", address, i))
		if err != nil {
			return make(map[int]int), err
		}
		defer response.Body.Close()

		//parse response
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return make(map[int]int), err
		}

		var output Output
		err = json.Unmarshal(body, &output)
		if err != nil {
			return make(map[int]int), err
		}

		routerMap[output.ReceiveFrom + 1] = i + 1
	}

	return routerMap, nil
}
