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

// These global variables hold build information. The Build variable will be
// populated by the Makefile and uses the Git Head hash as its identifier.
// These variables are used in the console output for --version and --help.
var (
	Version = "0.0.1"
	Build   string
)

func main() {
	configFileName := processCommandLineFlags()

	// --------------------------------------------------
	// Define variables
	// --------------------------------------------------

	router := mux.NewRouter()
	serviceCounter := 0
	var config config.ServerConfigType
	config.Router = router

	// --------------------------------------------------
	// Load System and Server Configuration
	// --------------------------------------------------

	// In addition to checking the configuration for completeness the verify
	// process will also populate some of the values.
	config.LoadServerConfig(configFileName)
	configError := config.VerifyServerConfig()
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
	if config.Logging.Enabled == true {
		logFile, err := os.OpenFile(config.Logging.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	// Start a Discovery Service handler
	// --------------------------------------------------
	// This will look to see if there are any Discovery services
	// defined in the config file. If there are, it will loop through the list and setup
	// handlers for each one of them. The HandleFunc passes in copy of the Discovery Resource
	// and the extra meta data that it needs to process the request.

	if config.DiscoveryServer.Enabled == true {
		for _, s := range config.DiscoveryServer.Services {
			if s.Enabled == true {

				var ts server.TAXIIServerHandlerType
				ts.NewDiscoveryHandler(s)
				ts.Resource = config.DiscoveryResources[s.ResourceID]

				log.Println("Starting TAXII Discovery service at:", s.ResourcePath)
				router.HandleFunc(s.ResourcePath, ts.TAXIIServerHandler).Methods("GET")
				serviceCounter++
			}
		}
	}

	// --------------------------------------------------
	// Start an API Root Service handler
	// Example: /api1/
	// --------------------------------------------------
	// This will look to see if there are any API Root services defined
	// in the config file. If there are, it will loop through the list and setup handlers
	// for each one of them. The HandleFunc passes in copy of the API Root Resource and the
	// extra meta data that it needs to process the request.

	if config.APIRootServer.Enabled == true {
		for _, api := range config.APIRootServer.Services {
			if api.Enabled == true {

				var ts server.TAXIIServerHandlerType
				ts.NewAPIRootHandler(api)
				ts.Resource = config.APIRootResources[api.ResourceID]

				log.Println("Starting TAXII API Root service at:", api.ResourcePath)
				router.HandleFunc(ts.ResourcePath, ts.TAXIIServerHandler).Methods("GET")
				serviceCounter++

				// --------------------------------------------------
				// Start a Collections Service handler
				// Example: /api1/collections/
				// --------------------------------------------------
				// This will look to see if the Collections service is enabled
				// in the configuration file for a given API Root. If it is, it
				// will setup handlers for it.
				// The HandleFunc passes in copy of the Collections Resource and the extra meta data
				// that it needs to process the request.

				if api.Collections.Enabled == true {

					var collectionsSrv server.TAXIIServerHandlerType
					collectionsSrv.NewCollectionsHandler(api)
					collections := objects.NewCollections()

					// We need to look in to this instance of the API Root and find out which collections are tied to it
					// Then we can use that ID to pull from the collections list and add them to this list of valid collections
					for _, c := range api.Collections.Members {

						// If enabled, only add the collection to the list if the collection can either be read or written to
						if config.CollectionResources[c].Resource.CanRead == true || config.CollectionResources[c].Resource.CanWrite == true {
							collections.AddCollection(config.CollectionResources[c].Resource)
						}
					}
					collectionsSrv.Resource = collections

					log.Println("Starting TAXII Collections service of:", api.Collections.ResourcePath)
					router.HandleFunc(collectionsSrv.ResourcePath, collectionsSrv.TAXIIServerHandler).Methods("GET")

					// --------------------------------------------------
					// Start a Collection handler
					// Example: /api1/collections/9cfa669c-ee94-4ece-afd2-f8edac37d8fd/
					// --------------------------------------------------
					// This will look to see which collections are defined for this
					// Collections group in this API Root. If they are enabled, it
					// will setup handlers for it.
					// The HandleFunc passes in copy of the Collection Resource and the extra meta data
					// that it needs to process the request.

					for _, c := range api.Collections.Members {

						resourceCollectionIDPath := collectionsSrv.ResourcePath + config.CollectionResources[c].Resource.ID + "/"

						// Make a copy of just the elements that we need to process the request and nothing more.
						// This is done to prevent sending the entire server config in to each handler
						var collectionSrv server.TAXIIServerHandlerType
						collectionSrv.NewCollectionHandler(api, resourceCollectionIDPath)
						collectionSrv.Resource = config.CollectionResources[c].Resource

						log.Println("Starting TAXII Collection service of:", resourceCollectionIDPath)

						// We do not need to check to see if the collection is enabled
						// and readable/writable because that was already done
						// TODO add support for post if the collection is writable
						router.HandleFunc(collectionSrv.ResourcePath, collectionSrv.TAXIIServerHandler).Methods("GET")

						// --------------------------------------------------
						// Start an Objects handler
						// Example: /api1/collections/9cfa669c-ee94-4ece-afd2-f8edac37d8fd/objects/
						// --------------------------------------------------

						// Make a copy of just the elements that we need to process the request and nothing more.
						// This is done to prevent sending the entire server config in to each handler
						var objectsSrv server.STIXServerHandlerType
						objectsSrv.Type = "Objects"
						objectsSrv.ResourcePath = resourceCollectionIDPath + "objects/"
						objectsSrv.HTMLEnabled = config.APIRootServer.HTMLEnabled
						objectsSrv.HTMLTemplateFile = api.HTMLBranding.Objects
						objectsSrv.HTMLTemplatePath = config.Global.Prefix + config.Global.HTMLTemplateDir
						objectsSrv.LogLevel = config.Logging.LogLevel

						// --------------------------------------------------
						// Start a Objects and Object by ID handlers
						// --------------------------------------------------

						log.Println("Starting TAXII Object service of:", objectsSrv.ResourcePath)
						config.Router.HandleFunc(objectsSrv.ResourcePath, objectsSrv.ObjectsServerHandler).Methods("GET")

						log.Println("Starting TAXII Object service of:", objectsSrv.ResourcePath)
						objectsSrv.ResourcePath = resourceCollectionIDPath + "objects/" + "{objectid}/"
						config.Router.HandleFunc(objectsSrv.ResourcePath, objectsSrv.ObjectsServerHandler).Methods("GET")
					} // End for loop api.Collections.Members
				} // End if Collections.Enabled == true
			}
		}
	}

	// --------------------------------------------------
	// Fail if no services are running
	// --------------------------------------------------
	if serviceCounter == 0 {
		log.Fatalln("No TAXII services defined")
	}

	// --------------------------------------------------
	// Listen for Incoming Connections
	// --------------------------------------------------
	if config.Global.Protocol == "http" {
		log.Println("Listening on:", config.Global.Listen)
		log.Fatalln(http.ListenAndServe(config.Global.Listen, router))
	} else if config.Global.Protocol == "https" {
		// --------------------------------------------------
		// Configure TLS settings
		// --------------------------------------------------
		// TODO move TLS elements to configuration file
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
			Addr:         config.Global.Listen,
			Handler:      router,
			TLSConfig:    tlsConfig,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		tlsKeyPath := "etc/tls/" + config.Global.TLSKey
		tlsCrtPath := "etc/tls/" + config.Global.TLSCrt
		log.Fatalln(tlsServer.ListenAndServeTLS(tlsCrtPath, tlsKeyPath))
	} else {
		log.Fatalln("No valid protocol was defined in the configuration file")
	} // end else
}

// --------------------------------------------------
// Private functions
// --------------------------------------------------

// processCommandLineFlags - This function will process the command line flags
// and will print the version or help information as needed.
func processCommandLineFlags() string {
	defaultServerConfigFilename := "etc/freetaxii.conf"
	sOptServerConfigFilename := getopt.StringLong("config", 'c', defaultServerConfigFilename, "System Configuration File", "string")
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
	return *sOptServerConfigFilename
}

// printOutputHeader - This function will print a header for all console output
func printOutputHeader() {
	fmt.Println("")
	fmt.Println("FreeTAXII Server")
	fmt.Println("Copyright: Bret Jordan")
	fmt.Println("Version:", Version)
	if Build != "" {
		fmt.Println("Build:", Build)
	}
	fmt.Println("")
}
