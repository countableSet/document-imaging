package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
)

var (
	configLocation   = ".config" + string(os.PathSeparator) + "document-imaging" + string(os.PathSeparator) + "scanner.json"
	configPath       = ".config" + string(os.PathSeparator) + "document-imaging"
	metadataLocation = ".config" + string(os.PathSeparator) + "document-imaging" + string(os.PathSeparator) + "metadata.json"
	conf             = &config{}
	meta             = &metadata{}
	scanImageCmd     = exec.Command("scanimage", "-L")
	scanImageRe      = regexp.MustCompile("^device `(?P<name>.*)'")
	scanImageFailRe  = regexp.MustCompile("^scanimage: open of device (.*) failed")
)

// Config struct to store scanner id information
// should look something like
// { ScannerID: "genesys:libusb:001:011" }
type config struct {
	ScannerID string `json:"ScannerID"`
}

// Metadata struct to store pdf metadata information
// should look something like
// { author: "rachel" }
type metadata struct {
	Author string `json:"author"`
}

// Fetches the scanner id from the config file, if the id is empty, fetch the id from the command and return that
func fetchScannerIDFromConfig() string {
	usr, _ := user.Current()
	configFile := usr.HomeDir + string(os.PathSeparator) + configLocation
	load(configFile, conf)
	configFile = usr.HomeDir + string(os.PathSeparator) + metadataLocation
	load(configFile, meta)
	return verifyConfig()
}

// Load the JSON config file from the provided path, and unmarshal into conf object
func load(configFile string, i interface{}) {
	var err error
	var input = io.ReadCloser(os.Stdin)
	if input, err = os.Open(configFile); err != nil {
		fmt.Println("warning: Could not find config file in " + configFile)
		return
	}

	// Read the config file
	jsonBytes, err := ioutil.ReadAll(input)
	input.Close()
	if err != nil {
		log.Fatalln(err)
	}

	// Parse the config
	if err := json.Unmarshal(jsonBytes, &i); err != nil {
		log.Fatalf("Could not parse %q: %v\n", configFile, err)
	}
}

// Verifies the config information, if the id is empty, it will fetch it and return that
// otherwise returns id from the config
func verifyConfig() string {
	var result string
	if conf.ScannerID == "" {
		result = fetchScannerID()
	} else {
		result = conf.ScannerID
	}
	return result
}

// Fetches the scanner id from the `scanimage -L` command, write it to the config file and returns it
// Command output: device `genesys:libusb:001:002' is a Canon LiDE 200 flatbed scanner
func fetchScannerID() string {
	output := runCommandWithOutput(scanImageCmd)
	if !scanImageRe.MatchString(output) {
		fmt.Fprintln(os.Stderr, "Cannot find scanner!")
		os.Exit(1)
	}
	matches := scanImageRe.FindAllStringSubmatch(output, -1)[0]
	match := matches[1]
	conf.ScannerID = match
	writeConfigToFile(configLocation, conf)
	return match
}

// Verifies the command output from scanning to determine if the scanner id is still valid
// Error output: scanimage: open of device genesys:libusb:001:011 failed: Invalid argument
func verifyScanCommandOutput(output string) (bool, string) {
	match := scanImageFailRe.FindString(output)
	if match == "" {
		return true, ""
	}
	fmt.Println("Device id has changed, fetching new one...")
	return false, fetchScannerID()
}

// Write the conf object to the config file in json format for reference
func writeConfigToFile(fileLocation string, i interface{}) {
	jsonOutput, err := json.Marshal(i)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot marshal config json!", err)
		os.Exit(1)
	}
	usr, _ := user.Current()
	path := usr.HomeDir + string(os.PathSeparator) + configPath
	err = os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error create path directories!", err)
		os.Exit(1)
	}
	configFile := usr.HomeDir + string(os.PathSeparator) + fileLocation
	err = ioutil.WriteFile(configFile, jsonOutput, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing to config json file!", err)
		os.Exit(1)
	}
}
