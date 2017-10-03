// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"errors"
	"log"
)

// verifyDisocveryConfig - This method will verify all of the configuration
// directives for the TAXII Discovery Service
func (ezt *ServerConfigType) verifyAPIRootConfig() error {
	var problemsFound = 0

	// This variable will track if any of the actual api root services are
	// enabled. If the outer service says yes, but no actual services are
	// enabled, throw an error.
	var isServiceEnabled = false

	// API Service Directives
	for i, value := range ezt.APIRootServer.Services {

		// Copy in logging level to make it easier to use in a handler
		ezt.APIRootServer.Services[i].LogLevel = ezt.Logging.LogLevel

		// If this service instance is enabled
		if value.Enabled == true {
			isServiceEnabled = true
		}

		// Verify the API Name is defined in the configuration file. This is used as the path name
		if value.Name == "" {
			log.Println("CONFIG: One or more API Root Services is missing the 'name' directive in the configuration file")
			problemsFound++
		} else {
			ezt.APIRootServer.Services[i].ResourcePath = "/" + value.Name + "/"
		}

		// Set the collections resource path
		ezt.APIRootServer.Services[i].Collections.ResourcePath = ezt.APIRootServer.Services[i].ResourcePath + "collections/"

		// Verify the API Resource is found
		if _, ok := ezt.APIRootResources[value.ResourceID]; !ok {
			value := "CONFIG: The API Root Resource " + value.ResourceID + " is missing from the configuration file"
			log.Println(value)
			problemsFound++
		}
	} // End for loop

	if isServiceEnabled == false {
		log.Println("CONFIG: While the API Root Server is enabled, there are no API Root Services that are enabled")
		problemsFound++
	}

	// Return errors if there were any
	if problemsFound == 1 {
		log.Println("ERROR:", problemsFound, "error was found in the API Root configuration")
		return errors.New("Configuration Errors Found")
	} else if problemsFound > 1 {
		log.Println("ERROR:", problemsFound, "errors were found in the API Root configuration")
		return errors.New("Configuration Errors Found")
	} else {
		return nil
	}
}
