package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ReadConfiguration(fname string) (map[string]map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("Failed to read configuration in ReadConfiguration(): %v", err)
	}

	var result map[string]map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal json in ReadConfiguration(): %v", err)
	}
	return result, nil
}
