package helpers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SwitchInput(address string, input string, output string) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/Port/Set/%s/%s", address, input, output))
	defer resp.Body.Close()

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("Pulse eight returned error code: %s and error %s", resp.StatusCode, responseBody))
	}

	return nil
}

