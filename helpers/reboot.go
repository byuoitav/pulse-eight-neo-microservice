package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/byuoitav/common/log"
)

// RebootNeo will run a command that will reboot any pulse 8
func RebootNeo(address string) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/System/Restart", address)) //TODO: Finish this command and get it once a neo is available to test on
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Pulse Eight returned error code: %v and error %s", resp.StatusCode, responseBody)
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
