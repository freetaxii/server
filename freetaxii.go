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
	"github.com/gorilla/mux"
	"github.com/pborman/getopt"
	"log"
	"net/http"
	"os"
)

// DEFAULT_SERVER_CONFIG_FILENAME - is a constant for the default location of the configuration file
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
	taxiiServerConfig.Router = router

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
	// Start a Discovery Service handler
	// --------------------------------------------------
	if taxiiServerConfig.DiscoveryService.Enabled == true {
		serviceCounter += taxiiServerConfig.StartDiscoveryService()
	}

	// --------------------------------------------------
	// Start an API Root Service handler
	// --------------------------------------------------
	if taxiiServerConfig.ApiRootService.Enabled == true {
		serviceCounter += taxiiServerConfig.StartApiRootService()
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
	if taxiiServerConfig.System.Protocol == "http" {
		log.Println("Listening on:", taxiiServerConfig.System.Listen)
		log.Fatalln(http.ListenAndServe(taxiiServerConfig.System.Listen, router))
	} else {
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
			Addr:         taxiiServerConfig.System.Listen,
			Handler:      router,
			TLSConfig:    tlsConfig,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		tlsKeyPath := "etc/tls/" + taxiiServerConfig.System.TlsKey
		tlsCrtPath := "etc/tls/" + taxiiServerConfig.System.TlsCrt
		log.Fatalln(tlsServer.ListenAndServeTLS(tlsCrtPath, tlsKeyPath))
	} // end else
}

// --------------------------------------------------
// Print Help and Version information
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
