package main

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestPromptUser(t *testing.T) {
	var tests = []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"y", true},
		{"yes", true},
		{"Y", true},
		{"YES", true},
		{"Yes", true},
		{"n", false},
		{"no", false},
		{"N", false},
		{"NO", false},
		{"No", false},
		{`t
y`, true},
		{`quit
y`, true},
	}
	for _, test := range tests {
		r := strings.NewReader(test.input)
		result := promptUser(r)
		if result != test.expected {
			t.Errorf("Expected %t but got %t", test.expected, result)
		}
	}
}

func TestPromptForFilename(t *testing.T) {
	var date = time.Now().Format("2006-01-02_")
	var testFiles = []string{date + "existing1.pdf", "2015-01-01existing2.pdf"}
	helperCreateExistingTestFiles(testFiles)
	defer helperRemoveExistingTestFiles(testFiles)
	var tests = []struct {
		input    string
		expected string
	}{
		{"unittest1", date + "unittest1"},
		{"2014-08-01_unittest1", "2014-08-01_unittest1"},
		{"2014-08-01_unittest1", "2014-08-01_unittest1"},
		{`existing1
unittest3`, date + "unittest3"},
		{`2015-01-01existing2
2015-01-01 unittest4`, "2015-01-01 unittest4"},
	}
	for _, test := range tests {
		r := strings.NewReader(test.input)
		result := promptForFilename(r)
		if result != test.expected {
			t.Errorf("Expected %s but got %s", test.expected, result)
		}
	}
}

func TestVerifyFilenameIncludesDate(t *testing.T) {
	var date = time.Now().Format("2006-01-02_")
	var tests = []struct {
		input    string
		expected string
	}{
		{"2014-01-01 newfile", "2014-01-01 newfile"},
		{"2014-01 newfile", "2014-01 newfile"},
		{"2014-01-01newfile", "2014-01-01newfile"},
		{"2014-01newfile", "2014-01newfile"},
		{"newfile", date + "newfile"},
	}
	for _, test := range tests {
		result := verifyFilenameIncludesDate(test.input)
		if result != test.expected {
			t.Errorf("Expected %s but got %s", test.expected, result)
		}
	}
}

func TestVerifyFileNotExist(t *testing.T) {
	var testFiles = []string{"existing1.pdf", "2015-01-01existing2.pdf"}
	helperCreateExistingTestFiles(testFiles)
	defer helperRemoveExistingTestFiles(testFiles)
	var tests = []struct {
		input    string
		expected bool
	}{
		{"existing1", false},
		{"2015-01-01existing2", false},
		{"untitest1", true},
	}
	for _, test := range tests {
		result := verifyFileNotExist(test.input)
		if result != test.expected {
			t.Errorf("Expected %t but got %t", test.expected, result)
		}
	}
}

func TestParseTitleAndDate(t *testing.T) {
	var tests = []struct {
		input         string
		expectedTitle string
		expectedDate  string
	}{
		{"2016-01-01_Test file", "Test file", "2016-01-01"},
		{"2016-01-01 Test file", "Test file", "2016-01-01"},
		{"2016-01-01-Test file", "Test file", "2016-01-01"},
		{"2016-01-01Test file", "Test file", "2016-01-01"},
		{"2016-01_Test file", "Test file", "2016-01"},
		{"2016-01-Test file", "Test file", "2016-01"},
		{"2016-01 Test file", "Test file", "2016-01"},
	}
	for _, test := range tests {
		resultTitle, resultDate := parseTitleAndDate(test.input)
		if resultTitle != test.expectedTitle {
			t.Errorf("Expected %s but got %s", test.expectedTitle, resultTitle)
		}
		if resultDate != test.expectedDate {
			t.Errorf("Expected %s but got %s", test.expectedDate, resultDate)
		}
	}
}

func TestNormalizeFilename(t *testing.T) {
	var tests = []struct {
		input            string
		expectedFilename string
	}{
		{"2018-01-01_unittest1", "2018-01-01_unittest1"},
		{"2018-01-01_unit test with spaces", "2018-01-01_unit_test_with_spaces"},
		{"2018-01_unittest2", "2018-01_unittest2"},
		{"2018-01 unit test with spaces2", "2018-01_unit_test_with_spaces2"},
		{"2018-01-01_UNITtest1", "2018-01-01_unittest1"},
		{"2018-01-01_UNit test WITH Spaces", "2018-01-01_unit_test_with_spaces"},
		{"2018-01_UNITtest2", "2018-01_unittest2"},
		{"2018-01 Unit Test With Spaces2", "2018-01_unit_test_with_spaces2"},
	}
	for _, test := range tests {
		resultFilename := normalizedFilename(test.input)
		if resultFilename != test.expectedFilename {
			t.Errorf("Expected %s but got %s", test.expectedFilename, resultFilename)
		}
	}
}

func helperCreateExistingTestFiles(files []string) {
	for _, file := range files {
		cmd := exec.Command("bash", "-c", "touch "+file)
		runCommand(cmd)
	}
}

func helperRemoveExistingTestFiles(files []string) {
	for _, file := range files {
		cmd := exec.Command("rm", "-f", file)
		runCommand(cmd)
	}
}
