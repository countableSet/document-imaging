package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// CommitHash contains the current Git revision. Use make to build to make
	// sure this gets set.
	CommitHash string

	// BuildDate contains the date of the current build.
	BuildDate string

	// Version Number contains vX.Y.Z of the current build
	VersionNumber string
)

var (
	version *bool
	author  *string
)

func flagParsing() {
	version = flag.Bool("v", false, "Print the version number")
	author = flag.String("a", "", "Configure metadata information for author and creator fields")
	flag.Parse()
	outputVersionInfo()
	saveMetadataInfo()
}

func outputVersionInfo() {
	if *version {
		fmt.Printf("Document Imaging %s Build Date: %s Commit Hash: %s\n", VersionNumber, BuildDate, CommitHash)
		os.Exit(0)
	}
}

func saveMetadataInfo() {
	if *author != "" {
		meta.Author = *author
		writeConfigToFile(metadataLocation, meta)
		os.Exit(0)
	}
}
