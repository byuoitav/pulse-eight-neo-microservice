package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/byuoitav/common/log"
)

type Result struct {
	Result bool `json:"Result"`
}

func SwitchInput(address string, input string, output string) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/Port/Set/%s/%s", address, input, output))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Pulse eight returned error code: %v and error %s", resp.StatusCode, responseBody)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.L.Infof("Response from %v: %s", address, body)

	var result Result
	json.Unmarshal(body, &result)

	if !result.Result {
		return fmt.Errorf("Pulse eight bad result: %v", result)
	}

	return nil
}
