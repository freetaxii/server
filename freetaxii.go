// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package main

import (
	"fmt"
	"github.com/freetaxii/freetaxii-server/lib/server"
	"github.com/gorilla/mux"
	"github.com/pborman/getopt"
	"log"
	"net/http"
	"os"
)

const (
	DEFAULT_SERVER_CONFIG_FILENAME = "etc/freetaxii.conf"
)

var sVersion = "0.0.1"

var sOptServerConfigFilename = getopt.StringLong("config", 'c', DEFAULT_SERVER_CONFIG_FILENAME, "System Configuration File", "string")
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
	var taxiiServer server.ServerType

	// --------------------------------------------------
	// Load System and Server Configuration
	// --------------------------------------------------

	taxiiServer.LoadServerConfig(*sOptServerConfigFilename)

	// --------------------------------------------------
	// Setup Logging File
	// --------------------------------------------------
	// TODO
	// Need to make the directory if it does not already exist
	// To do this, we need to split the filename from the directory, we will want to only
	// take the last bit in case there is multiple directories /etc/foo/bar/stuff.log

	// Only enable logging to a file if it is turned on in the configuration file
	if taxiiServer.Logging.Enabled == true {
		logFile, err := os.OpenFile(taxiiServer.Logging.LogFileFullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

	// This will look to see if there are any discovery services defined in the database
	// If there are, it will loop through the list and setup handlers for each one of them
	// The HandleFunc passes in an index value so that the handler instance will know
	// which Discovery Service it is processing. Without that information it can not
	// build the correct Discovery response message.
	if taxiiServer.DiscoveryService.Enabled == true {
		for i, _ := range taxiiServer.DiscoveryService.Resources {
			var index int = i

			// Check to see if this entry is actually enabled
			if taxiiServer.DiscoveryService.Resources[index].Enabled == true {
				var path string = taxiiServer.DiscoveryService.Resources[index].Path

				log.Println("Starting TAXII Discovery service at:", path)
				router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) { taxiiServer.DiscoveryServerHandler(w, r, index) }).Methods("GET")
				serviceCounter++
			}
		}
	}

	// // --------------------------------------------------
	// // Setup API Root Server
	// // --------------------------------------------------

	// if taxiiServer.ApiRootService.Enabled == true {
	// 	for i, _ := range taxiiServer.ApiRootService.Resources {
	// 		var index int = i
	// 		var path string = taxiiServer.ApiRootService.Resources[index].Path

	// 		log.Println("Starting TAXII API Root service at:", taxiiServer.ApiRootService.Resources[index].Path)
	// 		router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) { taxiiServer.ApiRootServerHandler(w, r, index) }).Methods("GET")
	// 		serviceCounter++
	// 	}
	// }

	// // --------------------------------------------------
	// // Setup Poll Server
	// // --------------------------------------------------

	// if taxiiServer.Services.Poll != "" {
	// 	log.Println("Starting TAXII Poll services at:", taxiiServer.Services.Poll)
	// 	http.HandleFunc(taxiiServer.Services.Poll, taxiiServer.PollServerHandler)
	// 	serviceCounter++
	// }

	// // --------------------------------------------------
	// // Setup Admin Server
	// // --------------------------------------------------

	// if taxiiServer.Services.Admin != "" {
	// 	log.Println("Starting TAXII Admin services at:", taxiiServer.Services.Admin)
	// 	http.HandleFunc(taxiiServer.Services.Admin, taxiiServer.AdminServerHandler)
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
	if taxiiServer.System.Listen != "" {
		log.Println("Listening on:", taxiiServer.System.Listen)
		http.ListenAndServe(taxiiServer.System.Listen, router)
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
