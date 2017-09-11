package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
		return errors.New(fmt.Sprintf("Pulse eight returned error code: %s and error %s", resp.StatusCode, responseBody))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Response received: %s", body)

	var result Result
	json.Unmarshal(body, &result)

	if !result.Result {
		return errors.New(fmt.Sprintf("Pulse eight bad result: %s", result))
	}

	return nil
}
