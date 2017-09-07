// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package main

import (
	"fmt"
	"github.com/freetaxii/freetaxii-server/lib/server"
	"github.com/freetaxii/libtaxii2/objects"
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
	var taxiiServerConfig server.ServerConfigType

	// --------------------------------------------------
	// Load System and Server Configuration
	// --------------------------------------------------

	taxiiServerConfig.LoadServerConfig(*sOptServerConfigFilename)

	// --------------------------------------------------
	// Setup Logging File
	// --------------------------------------------------
	// TODO
	// Need to make the directory if it does not already exist
	// To do this, we need to split the filename from the directory, we will want to only
	// take the last bit in case there is multiple directories /etc/foo/bar/stuff.log

	// Only enable logging to a file if it is turned on in the configuration file
	if taxiiServerConfig.Logging.Enabled == true {
		logFile, err := os.OpenFile(taxiiServerConfig.Logging.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	// Start a Discovery handler
	// --------------------------------------------------
	// This will look to see if there are any Discovery services defined in the config file.
	// If there are, it will loop through the list and setup handlers for each one of them
	// The HandleFunc passes in copy of the Discovery Resource and the extra meta data
	// that it needs to process the request.
	if taxiiServerConfig.DiscoveryService.Enabled == true {
		for i, _ := range taxiiServerConfig.DiscoveryService.Services {
			var index int = i

			// Check to see if this entry is actually enabled
			if taxiiServerConfig.DiscoveryService.Services[index].Enabled == true {

				// Make a copy of just the elements that we need to process the request and nothing more.
				// This is done to prevent sending the entire server config in to each handler
				var taxiiDiscovery server.ServerHandlerType
				taxiiDiscovery.Path = taxiiServerConfig.DiscoveryService.Services[index].Path
				taxiiDiscovery.HtmlDir = taxiiServerConfig.System.HtmlDir
				taxiiDiscovery.LogLevel = taxiiServerConfig.Logging.LogLevel
				taxiiDiscovery.Resource = taxiiServerConfig.DiscoveryService.Services[index].Resource

				log.Println("Starting TAXII Discovery service at:", taxiiDiscovery.Path)
				router.HandleFunc(taxiiDiscovery.Path, taxiiDiscovery.DiscoveryServerHandler).Methods("GET")
				serviceCounter++
			}
		}
	}

	// --------------------------------------------------
	// Start an API Root handler
	// --------------------------------------------------
	// This will look to see if there are any API Root services defined in the config file.
	// If there are, it will loop through the list and setup handlers for each one of them
	// The HandleFunc passes in copy of the API Root Resource and the extra meta data
	// that it needs to process the request.
	if taxiiServerConfig.ApiRootService.Enabled == true {
		for i, _ := range taxiiServerConfig.ApiRootService.Services {
			var index int = i

			// Check to see if this entry is actually enabled
			if taxiiServerConfig.ApiRootService.Services[index].Enabled == true {

				// Make a copy of just the elements that we need to process the request and nothing more.
				// This is done to prevent sending the entire server config in to each handler
				var taxiiApiRoot server.ServerHandlerType
				taxiiApiRoot.Path = taxiiServerConfig.ApiRootService.Services[index].Path
				taxiiApiRoot.HtmlDir = taxiiServerConfig.System.HtmlDir
				taxiiApiRoot.LogLevel = taxiiServerConfig.Logging.LogLevel
				taxiiApiRoot.Resource = taxiiServerConfig.ApiRootService.Services[index].Resource

				log.Println("Starting TAXII API Root service at:", taxiiApiRoot.Path)
				router.HandleFunc(taxiiApiRoot.Path, taxiiApiRoot.ApiRootServerHandler).Methods("GET")

				// --------------------------------------------------
				// Start a Collections handler
				// --------------------------------------------------
				log.Println("Starting TAXII Collections for API Root:", taxiiApiRoot.Path)

				// Make a copy of just the elements that we need to process the request and nothing more.
				// This is done to prevent sending the entire server config in to each handler
				var taxiiCollections server.ServerHandlerType
				taxiiCollections.Path = taxiiServerConfig.ApiRootService.Services[index].Path + "collections/"
				taxiiCollections.HtmlDir = taxiiServerConfig.System.HtmlDir
				taxiiCollections.LogLevel = taxiiServerConfig.Logging.LogLevel

				// We need to look in to this instance of the API Root and find out which collections are tied to it
				// Then we can use that ID to pull from the collections list and add them to this collection
				collections := objects.NewCollections()
				for _, value := range taxiiServerConfig.ApiRootService.Services[index].Collections {

					// Only add the collection if it is enabled
					if taxiiServerConfig.Collections[value].Enabled == true {

						// If enabled, only add the collection to the list if the collection can either be read or written to
						if taxiiServerConfig.Collections[value].Resource.Can_read == true || taxiiServerConfig.Collections[value].Resource.Can_write == true {
							collections.AddCollection(taxiiServerConfig.Collections[value].Resource)
						}
					}

				}
				taxiiCollections.Resource = collections

				router.HandleFunc(taxiiCollections.Path, taxiiCollections.CollectionsServerHandler).Methods("GET")
				serviceCounter++
			}
		}
	}

	// // --------------------------------------------------
	// // Setup Admin Server
	// // --------------------------------------------------

	// if taxiiServerConfig.Services.Admin != "" {
	// 	log.Println("Starting TAXII Admin services at:", taxiiServerConfig.Services.Admin)
	// 	http.HandleFunc(taxiiServerConfig.Services.Admin, taxiiServerConfig.AdminServerHandler)
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
	if taxiiServerConfig.System.Listen != "" {
		log.Println("Listening on:", taxiiServerConfig.System.Listen)
		http.ListenAndServe(taxiiServerConfig.System.Listen, router)
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
