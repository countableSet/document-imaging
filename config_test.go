package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"testing"
)

func TestLoad(t *testing.T) {
	location := "load-test-file.json"
	defer func() {
		os.Remove(location)
	}()
	load(location)
	if conf.ScannerId != "" {
		t.Errorf("Expected not nil but got %s", conf.ScannerId)
	}
	id := "test"
	testConf := config{ScannerId: id}
	jsonOutput, _ := json.Marshal(testConf)
	ioutil.WriteFile(location, jsonOutput, 0666)
	load(location)
	if conf.ScannerId != id {
		t.Errorf("Expected %s but got %s", id, conf.ScannerId)
	}
}

func TestWriteConfigToFile(t *testing.T) {
	usr, _ := user.Current()
	configLocation = ".config" + string(os.PathSeparator) + "document-imaging" + string(os.PathSeparator) + "scanner-testing.json"
	location := usr.HomeDir + string(os.PathSeparator) + configLocation
	defer func() {
		configLocation = ".config" + string(os.PathSeparator) + "document-imaging" + string(os.PathSeparator) + "scanner.json"
		os.Remove(location)
	}()
	conf.ScannerId = ""
	writeConfigToFile()
	if _, err := os.Stat(location); os.IsNotExist(err) {
		t.Errorf("Expected path to exist %s", location)
	}
	if conf.ScannerId != "" {
		t.Errorf("Expected empty string but got %s", conf.ScannerId)
	}
	load(location)
	if conf.ScannerId != "" {
		t.Errorf("Expected empty string but got %s", conf.ScannerId)
	}
	id := "test"
	conf.ScannerId = id
	writeConfigToFile()
	if _, err := os.Stat(location); os.IsNotExist(err) {
		t.Errorf("Expected path to exist %s", location)
	}
	if conf.ScannerId != id {
		t.Errorf("Expected %s but got %s", id, conf.ScannerId)
	}
	load(location)
	if conf.ScannerId != id {
		t.Errorf("Expected %s but got %s", id, conf.ScannerId)
	}
}
