// Copyright 2015-2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/freetaxii/libstix2/datastore"
	"github.com/freetaxii/libstix2/datastore/sqlite3"
	"github.com/freetaxii/libstix2/resources/collections"
	"github.com/freetaxii/server/internal/config"
	"github.com/freetaxii/server/internal/handlers"
	"github.com/gologme/log"
	"github.com/gorilla/mux"
	"github.com/pborman/getopt"
)

/*
These global variables hold build information. The Build variable will be
populated by the Makefile and uses the Git Head hash as its identifier.
These variables are used in the console output for --version and --help.
*/
var (
	Version = "0.3.1"
	Build   string
)

func main() {
	configFileName := processCommandLineFlags()

	// Keep track of the number of services that are started
	serviceCounter := 0

	// --------------------------------------------------
	// Setup logger
	// --------------------------------------------------
	logger := log.New(os.Stderr, "", log.LstdFlags)

	// --------------------------------------------------
	// Load System and Server Configuration
	// --------------------------------------------------
	config, configError := config.New(logger, configFileName)
	if configError != nil {
		logger.Fatalln(configError)
	}

	// --------------------------------------------------
	// Setup Logging Levels
	// --------------------------------------------------
	switch config.Logging.Level {
	case 1:
		logger.EnableLevel("info")
	case 3:
		logger.EnableLevel("info")
		logger.EnableLevel("warn")
	case 5:
		logger.EnableLevel("info")
		logger.EnableLevel("warn")
		logger.EnableLevel("debug")
	case 10:
		logger.EnableLevel("info")
		logger.EnableLevel("warn")
		logger.EnableLevel("debug")
		logger.EnableLevel("trace")
	}

	logger.Traceln("TRACE: System Configuration Dump")
	logger.Tracef("%+v\n", config)

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
			logger.Fatalf("ERROR: can not open file: %v", err)
		}
		defer logFile.Close()
		logger.SetOutput(logFile)
	}

	// --------------------------------------------------
	// Setup Database Connection
	// --------------------------------------------------
	var ds datastore.Datastorer
	switch config.Global.DbType {
	case "sqlite3":
		databaseFilename := config.Global.Prefix + config.Global.DbFile
		ds = sqlite3.New(logger, databaseFilename, config.CollectionResources)
	default:
		logger.Fatalln("ERROR: unknown database type, or no database type defined in the server global configuration")
	}
	defer ds.Close()

	// --------------------------------------------------
	//
	// Configure HTTP Router
	//
	// --------------------------------------------------

	router := mux.NewRouter()
	config.Router = router

	// --------------------------------------------------
	//
	// Start Server
	//
	// --------------------------------------------------

	logger.Println("Starting FreeTAXII Server Version:", Version)

	// --------------------------------------------------
	//
	// Start a Discovery Service handler
	//
	// --------------------------------------------------
	// This will look to see if there are any Discovery services defined in the
	// configuration file. If there are, loop through the list and setup handlers
	// for each one of them. The HandleFunc takes in a copy of the Discovery
	// Resource and the extra meta data that it needs to process the request.

	if config.DiscoveryServer.Enabled == true {
		for _, s := range config.DiscoveryServer.Services {
			if s.Enabled == true {

				// Configuration for this specific instance and its resource
				ts, _ := handlers.NewDiscoveryHandler(logger, s, config.DiscoveryResources[s.ResourceID])

				logger.Infoln("Starting TAXII GET Discovery service at:", s.Path)
				router.HandleFunc(s.Path, ts.DiscoveryHandler).Methods("GET")
				serviceCounter++
			}
		}
	}

	// --------------------------------------------------
	// Start an API Root Service handler
	// Example: /api1/
	// --------------------------------------------------
	// This will look to see if there are any API Root services defined
	// in the config file. If there are, it will loop through the list
	// and setup handlers for each one of them. The HandleFunc passes in
	// copy of the API Root Resource and the extra meta data that it
	// needs to process the request.

	if config.APIRootServer.Enabled == true {
		for _, api := range config.APIRootServer.Services {
			if api.Enabled == true {

				logger.Infoln("Starting TAXII GET API Root service at:", api.Path)
				ts, _ := handlers.NewAPIRootHandler(logger, api, config.APIRootResources[api.ResourceID])
				router.HandleFunc(api.Path, ts.APIRootHandler).Methods("GET")
				serviceCounter++

				// Loop through the collections, if enabled and start the endpoints
				if api.Collections.Enabled == true {
					// Make a new map so we can work on a copy, this way we can
					// keep permissions unique per API root.
					colResources := make(map[string]*collections.Collection)

					// For each collection listed with ReadAccess add it to our local
					// copy called colResources and set the CanRead to true
					for _, c := range api.Collections.ReadAccess {
						if _, found := colResources[c]; !found {
							a := config.CollectionResources[c]
							colResources[c] = &a
							colResources[c].CanRead = true
						}
					}

					// For each collection listed with WriteAccess add it to our
					// local copy, only if it is not already found and set the
					// CanWrite to true
					for _, c := range api.Collections.WriteAccess {
						if _, found := colResources[c]; !found {
							a := config.CollectionResources[c]
							colResources[c] = &a
						}
						colResources[c].CanWrite = true
					}

					// Loop through all of the possible collections that are part
					// of this API Root and have either CanRead or CanWrite access
					// and add them to the Collection. This will prevent any collections
					// from showing up in the list if they do not have at least
					// read or write permissions.
					collections := collections.New()
					for key, _ := range colResources {
						col := colResources[key]
						collections.AddCollection(col)
					}

					// --------------------------------------------------
					// Start a Collections Service handler
					// Example: /api1/collections/
					// --------------------------------------------------
					collectionsSrv, _ := handlers.NewCollectionsHandler(logger, api, *collections, config.Global.ServerRecordLimit)
					logger.Infoln("Starting TAXII GET Collections service of:", collectionsSrv.URLPath)
					router.HandleFunc(collectionsSrv.URLPath, collectionsSrv.CollectionsHandler).Methods("GET")

					// Loop through all the collections that we have identified
					// that should have basic read or write access.
					for _, collectionResourse := range colResources {

						// --------------------------------------------------
						// Start a Collection handler
						// Example: /api1/collections/9cfa669c-ee94-4ece-afd2-f8edac37d8fd/
						// --------------------------------------------------
						// We do not need to check to see if the collection is enabled because that was already done
						collectionSrv, _ := handlers.NewCollectionHandler(logger, api, *collectionResourse, config.Global.ServerRecordLimit)
						logger.Infoln("Starting TAXII GET Collection service of:", collectionSrv.URLPath)
						router.HandleFunc(collectionSrv.URLPath, collectionSrv.CollectionHandler).Methods("GET")

						// --------------------------------------------------
						// Start an Objects handler
						// Example: /api1/collections/9cfa669c-ee94-4ece-afd2-f8edac37d8fd/objects/
						// --------------------------------------------------
						srvObjects, _ := handlers.NewObjectsHandler(logger, api, collectionResourse.ID, config.Global.ServerRecordLimit)
						srvObjects.DS = ds

						if collectionResourse.CanRead == true {
							logger.Infoln("Starting TAXII GET Object service of:", srvObjects.URLPath)
							config.Router.HandleFunc(srvObjects.URLPath, srvObjects.STIXContentServerHandler).Methods("GET")
						}

						if collectionResourse.CanWrite == true {
							logger.Infoln("Starting TAXII POST Object service of:", srvObjects.URLPath)
							config.Router.HandleFunc(srvObjects.URLPath, srvObjects.ObjectsServerWriteHandler).Methods("POST")
						}

						// --------------------------------------------------
						// Start a Objects by ID handlers
						// Example: /api1/collections/9cfa669c-ee94-4ece-afd2-f8edac37d8fd/objects/{objectid}/
						// --------------------------------------------------
						srvObjectsByID, _ := handlers.NewObjectsByIDHandler(logger, api, collectionResourse.ID, config.Global.ServerRecordLimit)
						srvObjectsByID.DS = ds

						if collectionResourse.CanRead == true {
							logger.Infoln("Starting TAXII GET Object by ID service of:", srvObjectsByID.URLPath)
							config.Router.HandleFunc(srvObjectsByID.URLPath, srvObjectsByID.STIXContentServerHandler).Methods("GET")
						}

						// --------------------------------------------------
						// Start a Objects by ID Versions handlers
						// Example: /api1/collections/9cfa669c-ee94-4ece-afd2-f8edac37d8fd/objects/{objectid}/versions/
						// --------------------------------------------------
						srvObjectVersions, _ := handlers.NewObjectVersionsHandler(logger, api, collectionResourse.ID, config.Global.ServerRecordLimit)
						srvObjectVersions.DS = ds

						if collectionResourse.CanRead == true {
							logger.Infoln("Starting TAXII GET Object Versions service of:", srvObjectVersions.URLPath)
							config.Router.HandleFunc(srvObjectVersions.URLPath, srvObjectVersions.STIXContentServerHandler).Methods("GET")
						}

						// --------------------------------------------------
						// Start a Manifest handler
						// Example: /api1/collections/9cfa669c-ee94-4ece-afd2-f8edac37d8fd/manifest/
						// --------------------------------------------------
						srvManifest, _ := handlers.NewManifestHandler(logger, api, collectionResourse.ID, config.Global.ServerRecordLimit)
						srvManifest.DS = ds

						if collectionResourse.CanRead == true {
							logger.Infoln("Starting TAXII GET Manifest service of:", srvManifest.URLPath)
							config.Router.HandleFunc(srvManifest.URLPath, srvManifest.STIXContentServerHandler).Methods("GET")
						}

					} // End for loop api.Collections.ResourceIDs
				} // End if Collections.Enabled == true
			} // End if api.Enabled == true
		} // End for loop API Root Services
	} // End if APIRootServer.Enabled == true

	// --------------------------------------------------
	//
	// Fail if no services are running
	//
	// --------------------------------------------------

	if serviceCounter == 0 {
		logger.Fatalln("No TAXII services defined")
	}

	// --------------------------------------------------
	//
	// Listen for Incoming Connections
	//
	// --------------------------------------------------

	if config.Global.Protocol == "http" {
		logger.Infoln("Listening on:", config.Global.Listen)
		logger.Fatalln(http.ListenAndServe(config.Global.Listen, router))
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
		logger.Fatalln(tlsServer.ListenAndServeTLS(tlsCrtPath, tlsKeyPath))
	} else {
		logger.Fatalln("No valid protocol was defined in the configuration file")
	} // end if statement
}

// --------------------------------------------------
//
// Private functions
//
// --------------------------------------------------

/*
processCommandLineFlags - This function will process the command line flags
and will print the version or help information as needed.
*/
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

/*
printOutputHeader - This function will print a header for all console output
*/
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
