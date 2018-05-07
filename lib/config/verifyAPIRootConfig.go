// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package config

import (
	"log"
	"strings"
)

// verifyDisocveryConfig - This method will verify all of the configuration
// directives for the TAXII Discovery Service
func (config *ServerConfigType) verifyAPIRootConfig() int {
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

		// Verify the API Path is defined in the configuration file.
		if value.Path == "" {
			log.Println("CONFIG: One or more API Root Services is missing the 'path' directive in the configuration file")
			problemsFound++
		} else {
			if !strings.HasSuffix(value.Path, "/") {
				log.Println("CONFIG: The path in one or more API Roots is missing the ending slash '/'")
				problemsFound++
			}
			if !strings.HasPrefix(value.Path, "/") {
				log.Println("CONFIG: The path in one or more API Roots is missing the starting slash '/'")
				problemsFound++
			}
			config.APIRootServer.Services[i].FullPath = value.Path
		}

		// Set the collections resource path
		config.APIRootServer.Services[i].Collections.FullPath = config.APIRootServer.Services[i].FullPath + "collections/"

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
	}
	return problemsFound
}
