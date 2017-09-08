// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package main

import (
	"crypto/tls"
	"fmt"
	"github.com/freetaxii/freetaxii-server/lib/config"
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
	var taxiiServerConfig config.ServerConfigType

	// --------------------------------------------------
	// Load System and Server Configuration
	// --------------------------------------------------

	taxiiServerConfig.LoadServerConfig(*sOptServerConfigFilename)
	configError := taxiiServerConfig.VerifyServerConfig()
	if configError != nil {
		log.Fatalln(configError)
	}

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
				taxiiDiscovery.Type = "Discovery"
				taxiiDiscovery.Path = taxiiServerConfig.DiscoveryService.Services[index].Path
				taxiiDiscovery.HtmlResourceFile = "discoveryResource.html"
				taxiiDiscovery.HtmlResourcePath = taxiiServerConfig.System.HtmlDir + "/" + taxiiDiscovery.HtmlResourceFile
				taxiiDiscovery.LogLevel = taxiiServerConfig.Logging.LogLevel
				taxiiDiscovery.Resource = taxiiServerConfig.DiscoveryService.Services[index].Resource

				log.Println("Starting TAXII Discovery service at:", taxiiDiscovery.Path)
				router.HandleFunc(taxiiDiscovery.Path, taxiiDiscovery.TaxiiServerHandler).Methods("GET")
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
				taxiiApiRoot.Type = "API-Root"
				taxiiApiRoot.Path = taxiiServerConfig.ApiRootService.Services[index].Path
				taxiiApiRoot.HtmlResourceFile = "apirootResource.html"
				taxiiApiRoot.HtmlResourcePath = taxiiServerConfig.System.HtmlDir + "/" + taxiiApiRoot.HtmlResourceFile
				taxiiApiRoot.LogLevel = taxiiServerConfig.Logging.LogLevel
				taxiiApiRoot.Resource = taxiiServerConfig.ApiRootService.Services[index].Resource

				log.Println("Starting TAXII API Root service at:", taxiiApiRoot.Path)
				router.HandleFunc(taxiiApiRoot.Path, taxiiApiRoot.TaxiiServerHandler).Methods("GET")
				serviceCounter++

				// --------------------------------------------------
				// Start a Collections handler
				// --------------------------------------------------

				// Make a copy of just the elements that we need to process the request and nothing more.
				// This is done to prevent sending the entire server config in to each handler
				var taxiiCollections server.ServerHandlerType
				taxiiCollections.Type = "Collections"
				taxiiCollections.Path = taxiiServerConfig.ApiRootService.Services[index].Path + "collections/"
				taxiiCollections.HtmlResourceFile = "collectionsResource.html"
				taxiiCollections.HtmlResourcePath = taxiiServerConfig.System.HtmlDir + "/" + taxiiCollections.HtmlResourceFile
				taxiiCollections.LogLevel = taxiiServerConfig.Logging.LogLevel

				// We need to look in to this instance of the API Root and find out which collections are tied to it
				// Then we can use that ID to pull from the collections list and add them to this list of valid collections
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

				log.Println("Starting TAXII Collections service of:", taxiiCollections.Path)
				router.HandleFunc(taxiiCollections.Path, taxiiCollections.TaxiiServerHandler).Methods("GET")

				// --------------------------------------------------
				// Start a Collection handler
				// --------------------------------------------------

				// We need to loop through each collection in this API Root and setup a handler for it.
				for i, value := range collections.Collections {

					// Make a copy of just the elements that we need to process the request and nothing more.
					// This is done to prevent sending the entire server config in to each handler
					var taxiiCollection server.ServerHandlerType
					taxiiCollection.Type = "Collection"
					taxiiCollection.Path = taxiiCollections.Path + value.Id + "/"
					taxiiCollection.HtmlResourceFile = "collectionResource.html"
					taxiiCollection.HtmlResourcePath = taxiiServerConfig.System.HtmlDir + "/" + taxiiCollection.HtmlResourceFile
					taxiiCollection.LogLevel = taxiiServerConfig.Logging.LogLevel
					taxiiCollection.Resource = collections.Collections[i]

					// --------------------------------------------------
					// Start a Collection handler
					// --------------------------------------------------
					log.Println("Starting TAXII Collection service of:", taxiiCollection.Path)

					// We do not need to check to see if the collection is enabled and readable/writeable because that was already done
					// TODO add support for post if the colleciton is writeable
					router.HandleFunc(taxiiCollection.Path, taxiiCollection.TaxiiServerHandler).Methods("GET")

					// --------------------------------------------------
					// Start a Objects handler
					// --------------------------------------------------
					var taxiiObjects server.ServerHandlerType
					taxiiObjects.Type = "Objects"
					taxiiObjects.Path = taxiiCollection.Path + "objects/"
					taxiiObjects.HtmlResourceFile = "objectsResource.html"
					taxiiObjects.HtmlResourcePath = taxiiServerConfig.System.HtmlDir + "/" + taxiiObjects.HtmlResourceFile
					taxiiObjects.LogLevel = taxiiServerConfig.Logging.LogLevel

					log.Println("Starting TAXII Object service of:", taxiiObjects.Path)
					router.HandleFunc(taxiiObjects.Path, taxiiObjects.ObjectsServerHandler).Methods("GET")
				}
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
	if taxiiServerConfig.System.Protocol == "http" {
		log.Println("Listening on:", taxiiServerConfig.System.Listen)
		log.Fatalln(http.ListenAndServe(taxiiServerConfig.System.Listen, router))
	} else {
		// --------------------------------------------------
		// Configure TLS settings
		// --------------------------------------------------
		tlsConfig := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
		tlsServer := &http.Server{
			Addr:         taxiiServerConfig.System.Listen,
			Handler:      router,
			TLSConfig:    tlsConfig,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		tlsKeyPath := "etc/tls/" + taxiiServerConfig.System.TlsKey
		tlsCrtPath := "etc/tls/" + taxiiServerConfig.System.TlsCrt
		log.Fatalln(tlsServer.ListenAndServeTLS(tlsCrtPath, tlsKeyPath))
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
