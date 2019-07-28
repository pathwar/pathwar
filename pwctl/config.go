package main

import (
	"encoding/json"
	"io/ioutil"

	pwctlconfig "pathwar.land/pwctl/config"
)

func getConfig() (*pwctlconfig.Config, error) {
	// load config file
	configJSON, err := ioutil.ReadFile("/pwctl.json")
	if err != nil {
		return nil, err
	}
	var config pwctlconfig.Config
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
