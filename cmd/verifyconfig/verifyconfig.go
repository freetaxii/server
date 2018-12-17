// Copyright 2015-2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package main

import (
	"fmt"
	"os"

	"github.com/freetaxii/server/internal/config"
	"github.com/gologme/log"
	"github.com/pborman/getopt"
)

// These global variables hold build information. The Build variable will be
// populated by the Makefile and uses the Git Head hash as its identifier.
// These variables are used in the console output for --version and --help.
var (
	Version = "0.0.1"
	Build   string
)

// These global variables are for dealing with command line options
var (
	defaultServerConfigFilename = "../etc/freetaxii.conf"
	sOptServerConfigFilename    = getopt.StringLong("config", 'c', defaultServerConfigFilename, "System Configuration File", "string")
	bOptHelp                    = getopt.BoolLong("help", 0, "Help")
	bOptVer                     = getopt.BoolLong("version", 0, "Version")
)

func main() {
	processCommandLineFlags()

	logger := log.New(os.Stderr, "", log.LstdFlags)
	logger.EnableLevel("info")
	logger.EnableLevel("debug")

	// --------------------------------------------------
	// Define variables
	// --------------------------------------------------

	_, err := config.New(logger, *sOptServerConfigFilename)

	if err != nil {
		logger.Fatalln(err)
	}

	// --------------------------------------------------
	// Load System and Server Configuration
	// --------------------------------------------------

	logger.Println("No errors found")
}

// --------------------------------------------------
// Private functions
// --------------------------------------------------

// processCommandLineFlags - This function will process the command line flags
// and will print the version or help information as needed.
func processCommandLineFlags() {
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
}

// printOutputHeader - This function will print a header for all console output
func printOutputHeader() {
	fmt.Println("")
	fmt.Println("FreeTAXII Server")
	fmt.Println("Copyright, Bret Jordan")
	fmt.Println("Version:", Version)
	if Build != "" {
		fmt.Println("Build:", Build)
	}
	fmt.Println("")
}
