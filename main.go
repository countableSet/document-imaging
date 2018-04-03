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
	datere    = regexp.MustCompile("^([0-9]{4}(-[0-9]{2}){1,2})(.*)")
	titlere   = regexp.MustCompile("^(_|-|\\s|\\.)(.*)")
)

func main() {
	flagParsing()
	id := fetchScannerIDFromConfig()
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
	continueScanning := true
	pageCounter := 1
	text := "Scanning page %d\n"
	for continueScanning {
		fmt.Printf(text, pageCounter)

		filename := time.Now().Format("20060102150405") + ".tiff"

		scanimageArgs := []string{
			"scanimage", deviceName, "--mode Color",
			"--resolution 300 -x 215.9 -y 279.4", "--format=tiff", "-p", ">", filename,
		}
		cmd := exec.Command("bash", "-c", strings.Join(scanimageArgs, " "))
		err := runCommandWithReturnedError(cmd)
		if err != nil {
			result, newID := verifyScanCommandOutput(err.Error())
			if !result {
				removeFile(filename)
				deviceName = "--device-name '" + newID + "'"
				continue
			}
		}

		files = append(files, filename)
		continueScanning = promptUser(r)
		pageCounter++
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
		outputFilename := fmt.Sprintf("scan%03d.tiff", i+1)
		scanfiles = append(scanfiles, outputFilename)

		commandString := []string{
			"convert", file, "-shave 15x15 -bordercolor white -border 15",
			"-deskew 40%", "-background white", "-level 10%,70%,1",
			"-blur 2", "+dither", "+repage", "+matte", "-compress Group4",
			"-colorspace gray", "-format tiff", outputFilename,
		}
		cmd := exec.Command("bash", "-c", strings.Join(commandString, " "))
		runCommand(cmd)
	}
	fmt.Println("Converting Done")
}

func createFinalPdfDocument(filename string) {
	fmt.Println("Converting tiff scans to pdf")
	pdfName := normalizedFilename(filename) + ".pdf"
	scanfiles = append(scanfiles, "doc.tiff")
	cmd := exec.Command("tiffcp", scanfiles...)
	runCommand(cmd)
	title, _ := parseTitleAndDate(filename)
	pdfArgs := []string{"-j", "-t", title, "-k", "scan"}
	if meta.Author != "" {
		pdfArgs = append(pdfArgs, "-c", meta.Author, "-a", meta.Author)
	}
	pdfArgs = append(pdfArgs, "-o", pdfName, "doc.tiff")
	cmd = exec.Command("tiff2pdf", pdfArgs...)
	runCommand(cmd)
	fmt.Println("Convert Done")
}

func parseTitleAndDate(filename string) (string, string) {
	matches := datere.FindAllStringSubmatch(filename, -1)[0]
	title := matches[3]
	date := matches[1]
	if titlere.MatchString(title) {
		title = titlere.FindAllStringSubmatch(title, -1)[0][2]
	}
	return title, date
}

func normalizedFilename(filename string) string {
	lowered := strings.ToLower(filename)
	replaced := strings.Replace(lowered, " ", "_", -1)
	return replaced
}

func removeIntermediateFiles() {
	filesToRemove := append(files, scanfiles...)
	for _, file := range filesToRemove {
		removeFile(file)
	}
	fmt.Println("Cleanup Done")
}

func removeFile(f string) {
	cmd := exec.Command("rm", "-f", f)
	runCommand(cmd)
}
