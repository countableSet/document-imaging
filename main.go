package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	files     []string
	scanfiles []string
	yesre     = regexp.MustCompile("y(es)?")
	nore      = regexp.MustCompile("n(o)?")
	datere    = regexp.MustCompile("^([0-9]{4}(-[0-9]{2}){1,2}).*")
)

func main() {
	flagParsing()
	id := fetchScannerIdFromConfig()
	startScanning(id)
	filename := promptForFilename(os.Stdin)
	convertFilesIntoNewTiffFiles()
	createFinalPdfDocument(filename)
	removeIntermediateFiles()
}

func runCommand(cmd *exec.Cmd) {
	err := runCommandWithReturnedError(cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running command: ", err)
		os.Exit(1)
	}
}

func runCommandWithOutput(cmd *exec.Cmd) string {
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running command: ", err, ": ", stderr.String())
		os.Exit(1)
	}
	return stdout.String()
}

func runCommandWithReturnedError(cmd *exec.Cmd) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		err = errors.New(stderr.String())
	}
	return err
}

func startScanning(id string) {
	deviceName := "--device-name '" + id + "'"
	r := os.Stdin
	continue_scanning := true
	page_counter := 1
	text := "Scanning page %d\n"
	for continue_scanning {
		fmt.Printf(text, page_counter)

		filename := time.Now().Format("20060102150405") + ".tiff"
		files = append(files, filename)

		scanimage_args := []string{
			"scanimage", deviceName, "--mode Color",
			"--resolution 300 -x 215.9 -y 279.4", "--format=tiff", "-p", ">", filename,
		}
		cmd := exec.Command("bash", "-c", strings.Join(scanimage_args, " "))
		err := runCommandWithReturnedError(cmd)
		if err != nil {
			result, newId := verifyScanCommandOutput(err.Error())
			if !result {
				deviceName = "--device-name '" + newId + "'"
				continue
			}
		}

		continue_scanning = promptUser(r)
		page_counter++
	}
}

func promptUser(r io.Reader) bool {
	snr := bufio.NewScanner(r)
	for {
		fmt.Print("Scan more pages? [Y/n] ")
		snr.Scan()
		response := snr.Text()
		response = strings.ToLower(response)
		if response == "" || yesre.MatchString(response) {
			return true
		} else if nore.MatchString(response) {
			return false
		}
	}
	return false
}

func promptForFilename(r io.Reader) string {
	snr := bufio.NewScanner(r)
	var response string
	for {
		fmt.Print("Enter pdf filename, no extension: ")
		snr.Scan()
		response = snr.Text()
		if response != "" {
			response = strings.Trim(verifyFilenameIncludesDate(response), " ")
			if verifyFileNotExist(response) {
				break
			}
		}
	}
	return response
}

func verifyFilenameIncludesDate(filename string) string {
	result := filename
	if !datere.MatchString(filename) {
		result = time.Now().Format("2006-01-02_") + filename
		fmt.Printf("Filename did not include a date, the new filename is %s.pdf\n", result)
	}
	return result
}

func verifyFileNotExist(filename string) bool {
	if _, err := os.Stat(filename + ".pdf"); os.IsNotExist(err) {
		// file does not already exist
		return true
	}
	fmt.Printf("File already exists, with name %s.pdf pick another\n", filename)
	return false
}

func convertFilesIntoNewTiffFiles() {
	fmt.Println("Converting tiff files into scan files")
	for i, file := range files {
		output_filename := fmt.Sprintf("scan%03d.tiff", i+1)
		scanfiles = append(scanfiles, output_filename)

		command_string := []string{
			"convert", file, "-deskew 40%", "-background white", "-level 10%,70%,1",
			"-blur 2", "+dither", "+repage", "+matte", "-compress Group4", "-colorspace gray",
			"-format tiff", output_filename,
		}
		cmd := exec.Command("bash", "-c", strings.Join(command_string, " "))
		runCommand(cmd)
	}
	fmt.Println("Converting Done")
}

func createFinalPdfDocument(filename string) {
	fmt.Println("Converting tiff scans to pdf")
	pdfName := filename + ".pdf"
	scanfiles = append(scanfiles, "doc.tiff")
	cmd := exec.Command("tiffcp", scanfiles...)
	runCommand(cmd)
	cmd = exec.Command("tiff2pdf", "-j", "-o", pdfName, "doc.tiff")
	runCommand(cmd)
	fmt.Println("Convert Done")
}

func removeIntermediateFiles() {
	files_to_remove := append(files, scanfiles...)
	for _, file := range files_to_remove {
		cmd := exec.Command("rm", "-f", file)
		runCommand(cmd)
	}
	fmt.Println("Cleanup Done")
}
