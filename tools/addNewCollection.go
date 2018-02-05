// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package main

import (
	"bufio"
	"fmt"
	"github.com/freetaxii/libstix2/datastore/sqlite3"
	"github.com/freetaxii/libstix2/objects"
	"github.com/pborman/getopt"
	"log"
	"os"
	"strings"
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
	defaultDatabaseFilename = "freetaxii.db"
	sOptDatabaseFilename    = getopt.StringLong("filename", 'f', defaultDatabaseFilename, "Database Filename", "string")
	bOptHelp                = getopt.BoolLong("help", 0, "Help")
	bOptVer                 = getopt.BoolLong("version", 0, "Version")
)

func main() {
	processCommandLineFlags()
	ds := sqlite3.New(databaseFilename)
	c := objects.NewCollection()
	c.NewID()

	reader := bufio.NewReader(os.Stdin)

	// Ask for title
	fmt.Print("Collection Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSuffix(title, "\n")
	if title != "" {
		c.SetTitle(title)
	} else {
		log.Fatalln("ERROR: The title can not be blank")
	}

	// Ask for descritpion
	fmt.Print("Description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSuffix(description, "\n")

	if description != "" {
		c.SetDescription(description)
	}

	fmt.Print("Can Read (y/n): ")
	canRead, _ := reader.ReadString('\n')
	canRead = strings.TrimSuffix(canRead, "\n")
	if canRead == "y" {
		c.SetCanRead()
	}

	fmt.Print("Can Write (y/n): ")
	canWrite, _ := reader.ReadString('\n')
	canWrite = strings.TrimSuffix(canWrite, "\n")
	if canWrite == "y" {
		c.SetCanWrite()
	}

	c.AddMediaType("application/vnd.oasis.stix+json")

	ds.AddTAXIIObject(c)
	ds.Close()
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
	fmt.Println("FreeTAXII - TAXII Table Creator")
	fmt.Println("Copyright: Bret Jordan")
	fmt.Println("Version:", Version)
	if Build != "" {
		fmt.Println("Build:", Build)
	}
	fmt.Println("")
}
