// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/freetaxii/libstix2/datastore"
	"github.com/freetaxii/libstix2/datastore/sqlite3"
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
	sOptCollectionID        = getopt.StringLong("collectionid", 'c', "", "Collection ID", "string")
	sOptSTIXID              = getopt.StringLong("stixid", 's', "", "Object ID", "string")
	sOptSTIXType            = getopt.StringLong("type", 't', "", "Object Type", "string")
	sOptVersion             = getopt.StringLong("stixversion", 'v', "", "Version", "string")
	sOptAddedAfter          = getopt.StringLong("date", 'd', "", "Added After", "string")
	sOptRangeBegin          = getopt.IntLong("begin", 'b', 0, "Range Begin", "int")
	sOptRangeEnd            = getopt.IntLong("end", 'e', 0, "Range End", "int")
	bOptHelp                = getopt.BoolLong("help", 0, "Help")
	bOptVer                 = getopt.BoolLong("version", 0, "Version")
)

func main() {
	processCommandLineFlags()
	var q datastore.QueryType

	q.CollectionID = *sOptCollectionID

	if *sOptSTIXID != "" {
		q.STIXID = strings.Split(*sOptSTIXID, ",")
	}

	if *sOptSTIXType != "" {
		q.STIXType = strings.Split(*sOptSTIXType, ",")
	}

	if *sOptVersion != "" {
		q.STIXVersion = strings.Split(*sOptVersion, ",")
	}

	if *sOptAddedAfter != "" {
		q.AddedAfter = strings.Split(*sOptAddedAfter, ",")
	}

	q.RangeBegin = *sOptRangeBegin
	q.RangeEnd = *sOptRangeEnd
	q.RangeMax = 5

	ds := sqlite3.New(*sOptDatabaseFilename)
	o, meta, err := ds.GetObjectsFromCollection(q)
	if err != nil {
		log.Fatalln(err)
	}

	var data []byte
	data, _ = json.MarshalIndent(*o, "", "    ")
	fmt.Println(string(data))

	fmt.Println("==========================================================================================")
	fmt.Println("Total Records: ", meta.Size)
	fmt.Println("X Header Date Added First: ", meta.DateAddedFirst)
	fmt.Println("X Header Date Added Last:  ", meta.DateAddedLast)
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

	if *sOptCollectionID == "" {
		log.Fatalln("Collection ID must not be empty.")
	}

}

// printOutputHeader - This function will print a header for all console output
func printOutputHeader() {
	fmt.Println("")
	fmt.Println("FreeTAXII - TAXII Get Collections")
	fmt.Println("Copyright: Bret Jordan")
	fmt.Println("Version:", Version)
	if Build != "" {
		fmt.Println("Build:", Build)
	}
	fmt.Println("")
}
