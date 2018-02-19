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
func (config *ServerConfigType) verifyAPIRootConfig() error {
	var problemsFound = 0

	// This variable will track if any of the actual api root services are
	// enabled. If the outer service says yes, but no actual services are
	// enabled, throw an error.
	var isServiceEnabled = false

	// API Service Directives
	for i, value := range config.APIRootServer.Services {

		// If this service instance is enabled
		if value.Enabled == true {
			isServiceEnabled = true
		}

		// Verify the API Name is defined in the configuration file. This is used as the path name
		if value.Name == "" {
			log.Println("CONFIG: One or more API Root Services is missing the 'name' directive in the configuration file")
			problemsFound++
		} else {
			config.APIRootServer.Services[i].ResourcePath = "/" + value.Name + "/"
		}

		// Set the collections resource path
		config.APIRootServer.Services[i].Collections.ResourcePath = config.APIRootServer.Services[i].ResourcePath + "collections/"

		// Verify the API Resource is found
		if _, ok := config.APIRootResources[value.ResourceID]; !ok {
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
	if problemsFound > 0 {
		log.Println("ERROR: The API Root configuration has", problemsFound, "error(s)")
		return errors.New("ERROR: Configuration errors found")
	}
	return nil
}
