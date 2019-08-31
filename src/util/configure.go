package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Updatedate string `json:"enddate"`
	Likedir    string `json:"likedir"`
	Wpdir      string `json:"wpdir"`
}

const CONFIG_FILENAME string = "config.json"

func ReadConfiguration() (Configuration, error) {
	var Conf Configuration
	byteValue, err := ioutil.ReadFile(CONFIG_FILENAME)
	if err != nil {
		return Configuration{}, err
	}
	err = json.Unmarshal(byteValue, &Conf)
	if err != nil {
		return Configuration{}, err
	} else {
		return Conf, nil
	}
}

func WriteConfiguration(conf Configuration) error {
	jsonBytes, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(CONFIG_FILENAME, jsonBytes, os.ModePerm)
}
