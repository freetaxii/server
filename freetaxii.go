// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package main

import (
	"fmt"
	"github.com/freetaxii/freetaxii-server/lib/server"
	"github.com/freetaxii/freetaxii-server/lib/systemconfig"
	"github.com/gorilla/mux"
	"github.com/pborman/getopt"
	"log"
	"net/http"
	"os"
)

const (
	DEFAULT_SYSTEM_CONFIG_FILENAME = "etc/freetaxii.conf"
)

var sVersion = "0.0.1"

var sOptSystemConfigFilename = getopt.StringLong("config", 'c', DEFAULT_SYSTEM_CONFIG_FILENAME, "System Configuration File", "string")
var bOptHelp = getopt.BoolLong("help", 0, "Help")
var bOptVer = getopt.BoolLong("version", 0, "Version")

func main() {
	getopt.HelpColumn = 35
	getopt.DisplayWidth = 120
	getopt.SetParameters("")
	getopt.Parse()

	if *bOptVer {
		printVersion()
	}

	if *bOptHelp {
		printHelp()
	}

	// --------------------------------------------------
	// Define variables
	// --------------------------------------------------

	router := mux.NewRouter()
	serviceCounter := 0
	var syscfg systemconfig.SystemConfigType
	var taxiiServer server.ServerType

	// --------------------------------------------------
	// Load System and Server Configuration
	// --------------------------------------------------

	syscfg.LoadSystemConfig(*sOptSystemConfigFilename)
	taxiiServer.SysConfig = &syscfg
	taxiiServer.LoadDiscoveryServicesConfig()

	// --------------------------------------------------
	// Setup Logging File
	// --------------------------------------------------
	// TODO
	// Need to make the directory if it does not already exist
	// To do this, we need to split the filename from the directory, we will want to only
	// take the last bit in case there is multiple directories /etc/foo/bar/stuff.log

	// Only enable logging to a file if it is turned on in the configuration file
	if syscfg.Logging.Enabled == true {
		logFile, err := os.OpenFile(syscfg.Logging.LogFileFullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer logFile.Close()

		log.SetOutput(logFile)
	}

	// --------------------------------------------------
	// Start Server
	// --------------------------------------------------
	log.Println("Starting FreeTAXII Server")

	// --------------------------------------------------
	// Setup Discovery Server
	// --------------------------------------------------

	// TODO set this up in a loop
	// How do you know which instance is being used? Without that you can not display the correct information from the object
	if taxiiServer.DiscoveryService.Enabled == true {
		for i, _ := range taxiiServer.DiscoveryService.Resources {
			var index int = i
			var path string = taxiiServer.DiscoveryService.Resources[index].Location

			log.Println("Starting TAXII Discovery services at:", taxiiServer.DiscoveryService.Resources[index].Location)
			router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) { taxiiServer.DiscoveryServerHandler(w, r, index) }).Methods("GET")
			serviceCounter++
		}
	}

	// // --------------------------------------------------
	// // Setup Collection Server
	// // --------------------------------------------------

	// if syscfg.Services.Collection != "" {
	// 	log.Println("Starting TAXII Collection services at:", syscfg.Services.Collection)
	// 	http.HandleFunc(syscfg.Services.Collection, taxiiServer.CollectionServerHandler)
	// 	serviceCounter++
	// }

	// // --------------------------------------------------
	// // Setup Poll Server
	// // --------------------------------------------------

	// if syscfg.Services.Poll != "" {
	// 	log.Println("Starting TAXII Poll services at:", syscfg.Services.Poll)
	// 	http.HandleFunc(syscfg.Services.Poll, taxiiServer.PollServerHandler)
	// 	serviceCounter++
	// }

	// // --------------------------------------------------
	// // Setup Admin Server
	// // --------------------------------------------------

	// if syscfg.Services.Admin != "" {
	// 	log.Println("Starting TAXII Admin services at:", syscfg.Services.Admin)
	// 	http.HandleFunc(syscfg.Services.Admin, taxiiServer.AdminServerHandler)
	// 	//serviceCounter++  Do not count this service in the list
	// }

	// --------------------------------------------------
	// Fail if no services are running
	// --------------------------------------------------

	if serviceCounter == 0 {
		log.Fatalln("No TAXII services defined")
	}

	// --------------------------------------------------
	// Listen for Incoming Connections
	// --------------------------------------------------

	// TODO - Need to verify the list address is a valid IPv4 address and port combination.
	if syscfg.System.Listen != "" {
		log.Println("Listening on:", syscfg.System.Listen)
		http.ListenAndServe(syscfg.System.Listen, router)
	} else {
		log.Fatalln("The listen directive is missing from the configuration file")
	}

}

// --------------------------------------------------
// Print Help and Version infomration
// --------------------------------------------------

func printHelp() {
	printOutputHeader()
	getopt.Usage()
	os.Exit(0)
}

func printVersion() {
	printOutputHeader()
	os.Exit(0)
}

// --------------------------------------------------
// Print a header for all output
// --------------------------------------------------

func printOutputHeader() {
	fmt.Println("")
	fmt.Println("FreeTAXII Server")
	fmt.Println("Copyright, Bret Jordan")
	fmt.Println("Version:", sVersion)
	fmt.Println("")
}
