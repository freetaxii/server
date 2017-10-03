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
func (ezt *ServerConfigType) verifyDiscoveryConfig() error {
	var problemsFound = 0

	// This variable will track if any of the actual discovery services are
	// enabled. If the outer service says yes, but no actual services are
	// enabled, throw an error.
	var isServiceEnabled = false

	// Loop through each Discovery Service and verify its configuration
	for i, value := range ezt.DiscoveryServer.Services {

		// Copy in logging level to make it easier to use in a handler
		ezt.DiscoveryServer.Services[i].LogLevel = ezt.Logging.LogLevel

		// Check to see if this service instance is enabled.
		if value.Enabled == true {
			isServiceEnabled = true
		}

		// Verify the Discovery Name is defined in the configuration file. This is used as the path name
		if value.Name == "" {
			log.Println("CONFIG: One or more Discovery Services is missing the 'name' directive in the configuration file")
			problemsFound++
		} else {
			ezt.DiscoveryServer.Services[i].ResourcePath = "/" + value.Name + "/"
		}

		// Verify the Discovery Resource that is referenced actually exists
		if _, ok := ezt.DiscoveryResources[value.ResourceID]; !ok {
			msg := "CONFIG: The Discovery Resource " + value.ResourceID + " is missing from the configuration file"
			log.Println(msg)
			problemsFound++
		}
	} // end for loop

	if isServiceEnabled == false {
		log.Println("CONFIG: While the Discovery Server is enabled, there are no Discovery Services that are enabled")
		problemsFound++
	}

	// Return errors if there were any
	if problemsFound > 0 {
		log.Println("ERROR: The Discovery configuration has", problemsFound, "error(s)")
		return errors.New("ERROR: Configuration errors found")
	}
	return nil
}
