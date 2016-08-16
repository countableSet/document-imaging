package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"os/exec"
	"regexp"
	"fmt"
)

var (
	configLocation = ".config"+string(os.PathSeparator)+"document-imaging"+string(os.PathSeparator)+"scanner.json"
	configPath = ".config"+string(os.PathSeparator)+"document-imaging"
	conf = &config{}
	scanImageCmd = exec.Command("scanimage", "-L")
	scanImageRe = regexp.MustCompile("^device `(?P<name>.*)'")
	scanImageFailRe = regexp.MustCompile("^scanimage: open of device (.*) failed")
)

// Config struct to store scanner id information
// should look something like
// {
//    scannerId: "genesys:libusb:001:011"
// }
type config struct {
	ScannerId string `json:"scannerId"`
}

func fetchScannerIdFromConfig() string {
	load()
	return verifyConfig()
}

// Load the JSON config file
func load() {
	usr, _ := user.Current()
	configFile := usr.HomeDir+string(os.PathSeparator)+configLocation

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
	if err := json.Unmarshal(jsonBytes, &conf); err != nil {
		log.Fatalln("Could not parse %q: %v", configFile, err)
	}
}

// device `genesys:libusb:001:002' is a Canon LiDE 200 flatbed scanner
func verifyConfig() string {
	var result string
	if conf.ScannerId == "" {
		result = fetchScannerId()
	} else {
		result = conf.ScannerId
	}
	return result
}

func fetchScannerId() string {
	output := runCommandWithOutput(scanImageCmd)
	if !scanImageRe.MatchString(output) {
		fmt.Fprintln(os.Stderr, "Cannot find scanner!")
		os.Exit(1)
	}
	matches := scanImageRe.FindAllStringSubmatch(output, -1)[0]
	match := matches[1]
	conf.ScannerId = match
	writeConfigToFile()
	return match
}

// scanimage: open of device genesys:libusb:001:011 failed: Invalid argument
func verifyScanCommandOutput(output string) (bool, string) {
	match := scanImageFailRe.FindString(output)
	if match == "" {
		return true, ""
	} else {
		return false, fetchScannerId()
	}
}

func writeConfigToFile() {
	fmt.Println(conf.ScannerId)
	jsonOutput, err := json.Marshal(conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot marshal config json!", err)
		os.Exit(1)
	}
	usr, _ := user.Current()
	path := usr.HomeDir+string(os.PathSeparator)+configPath
	err = os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error create path directories!", err)
		os.Exit(1)
	}
	configFile := usr.HomeDir+string(os.PathSeparator)+configLocation
	err = ioutil.WriteFile(configFile, jsonOutput, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing to config json file!", err)
		os.Exit(1)
	}
}
