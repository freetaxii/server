// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package main

import (
	"fmt"
	"github.com/freetaxii/libstix2/datastore/sqlite3"
	"github.com/pborman/getopt"
	"os"
)

// These global variables hold build information. The Build variable will be
// populated by the Makefile and uses the Git Head hash as its identifier.
// These variables are used in the console output for --version and --help.
var (
	Version = "0.0.1"
	Build   string
)

func main() {
	databaseFilename := processCommandLineFlags()
	ds := sqlite3.New(databaseFilename)

	ds.CreateAllSTIXTables()
	ds.CreateAllVocabTables()
	ds.PopulateAllVocabTables()
	ds.CreateAllTAXIITables()

	ds.Close()
}

// --------------------------------------------------
// Private functions
// --------------------------------------------------

// processCommandLineFlags - This function will process the command line flags
// and will print the version or help information as needed.
func processCommandLineFlags() string {
	defaultDatabaseFilename := "freetaxii.db"
	sOptDatabaseFilename := getopt.StringLong("filename", 'f', defaultDatabaseFilename, "Database Filename", "string")
	bOptHelp := getopt.BoolLong("help", 0, "Help")
	bOptVer := getopt.BoolLong("version", 0, "Version")

	getopt.HelpColumn = 35
	getopt.DisplayWidth = 120
	getopt.SetParameters("")
	getopt.Parse()

	// Lets check to see if the version command line flag was given. If it is
	// lets print out the version infomration and exit.
	if *bOptVer {
		printOutputHeader()
		os.Exit(0)
	}

	// Lets check to see if the help command line flag was given. If it is lets
	// print out the help information and exit.
	if *bOptHelp {
		printOutputHeader()
		getopt.Usage()
		os.Exit(0)
	}
	return *sOptDatabaseFilename
}

// printOutputHeader - This function will print a header for all console output
func printOutputHeader() {
	fmt.Println("")
	fmt.Println("FreeTAXII - STIX Table Creator")
	fmt.Println("Copyright: Bret Jordan")
	fmt.Println("Version:", Version)
	if Build != "" {
		fmt.Println("Build:", Build)
	}
	fmt.Println("")
}
