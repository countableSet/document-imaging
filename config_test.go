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
	load(location, conf)
	if conf.ScannerID != "" {
		t.Errorf("Expected not nil but got %s", conf.ScannerID)
	}
	id := "test"
	testConf := config{ScannerID: id}
	jsonOutput, _ := json.Marshal(testConf)
	ioutil.WriteFile(location, jsonOutput, 0666)
	load(location, conf)
	if conf.ScannerID != id {
		t.Errorf("Expected %s but got %s", id, conf.ScannerID)
	}
}

func TestWriteConfigToFile(t *testing.T) {
	usr, _ := user.Current()
	configLocation = "scanner-testing.json"
	location := usr.HomeDir + homeDirSubPath + string(os.PathSeparator) + configLocation
	defer func() {
		configLocation = ".config" + string(os.PathSeparator) + "document-imaging" + string(os.PathSeparator) + "scanner.json"
		os.Remove(location)
	}()
	conf.ScannerID = ""
	writeConfigToFile(configLocation, conf)
	if _, err := os.Stat(location); os.IsNotExist(err) {
		t.Errorf("Expected path to exist %s", location)
	}
	if conf.ScannerID != "" {
		t.Errorf("Expected empty string but got %s", conf.ScannerID)
	}
	load(location, conf)
	if conf.ScannerID != "" {
		t.Errorf("Expected empty string but got %s", conf.ScannerID)
	}
	id := "test"
	conf.ScannerID = id
	writeConfigToFile(configLocation, conf)
	if _, err := os.Stat(location); os.IsNotExist(err) {
		t.Errorf("Expected path to exist %s", location)
	}
	if conf.ScannerID != id {
		t.Errorf("Expected %s but got %s", id, conf.ScannerID)
	}
	load(location, conf)
	if conf.ScannerID != id {
		t.Errorf("Expected %s but got %s", id, conf.ScannerID)
	}
}
