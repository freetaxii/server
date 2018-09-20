// Copyright 2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package config

import (
	"strings"
)

/*
verifyDisocveryConfig - This method will verify all of the configuration
directives for the TAXII Discovery Service
*/
func (c *ServerConfig) verifyDiscoveryConfig() int {
	var problemsFound = 0

	// This variable will track if any of the actual discovery services are
	// enabled. If the outer service says yes, but no actual services are
	// enabled, throw an error.
	var isServiceEnabled = false

	// Loop through each Discovery Service and verify its configuration
	for i, value := range c.DiscoveryServer.Services {

		// Check to see if this service instance is enabled.
		if value.Enabled == true {
			isServiceEnabled = true
		}

		// Verify the Discovery Path is defined in the configuration file and set
		// the "full path" value
		if value.Path == "" {
			c.Logger.Println("CONFIG: One or more Discovery Services is missing the 'path' directive in the configuration file")
			problemsFound++
		} else {
			if !strings.HasSuffix(value.Path, "/") {
				c.Logger.Println("CONFIG: The path in one or more Discovery Services is missing the ending slash '/'")
				problemsFound++
			}
			if !strings.HasPrefix(value.Path, "/") {
				c.Logger.Println("CONFIG: The path in one or more Discovery Services is missing the starting slash '/'")
				problemsFound++
			}
		}

		// Verify the Discovery Resource that is referenced actually exists
		if _, ok := c.DiscoveryResources[value.ResourceID]; !ok {
			msg := "CONFIG: The Discovery Resource " + value.ResourceID + " is missing from the configuration file"
			c.Logger.Println(msg)
			problemsFound++
		}
	} // end for loop

	// Log an error if there are no Discovery services actually enabled.
	if isServiceEnabled == false {
		c.Logger.Println("CONFIG: While the Discovery Server is enabled, there are no Discovery Services that are enabled")
		problemsFound++
	}

	// Return errors if there were any
	if problemsFound > 0 {
		c.Logger.Println("ERROR: The Discovery configuration has", problemsFound, "error(s)")
	}
	return problemsFound
}

/*
verifyAPIRootConfig - This method will verify all of the configuration
directives for the TAXII API Root Service
*/
func (c *ServerConfig) verifyAPIRootConfig() int {
	var problemsFound = 0

	// This variable will track if any of the actual api root services are
	// enabled. If the outer service says yes, but no actual services are
	// enabled, throw an error.
	var isServiceEnabled = false

	// API Service Directives
	for i, value := range c.APIRootServer.Services {

		// If this service instance is enabled
		if value.Enabled == true {
			isServiceEnabled = true
		}

		// Verify the API Path is correctly defined in the configuration file.
		// Example: "path": "/api1/"
		if value.Path == "" {
			c.Logger.Println("CONFIG: One or more API Root Services is missing the 'path' directive in the configuration file")
			problemsFound++
		} else {
			if !strings.HasSuffix(value.Path, "/") {
				c.Logger.Println("CONFIG: The path in one or more API Roots is missing the ending slash '/'")
				problemsFound++
			}
			if !strings.HasPrefix(value.Path, "/") {
				c.Logger.Println("CONFIG: The path in one or more API Roots is missing the starting slash '/'")
				problemsFound++
			}
		}

		// Verify the API Resource is found
		if _, ok := c.APIRootResources[value.ResourceID]; !ok {
			value := "CONFIG: The API Root Resource " + value.ResourceID + " is missing from the configuration file"
			c.Logger.Println(value)
			problemsFound++
		}

		// Verify the Collection Resources are found
		if value.Collections.Enabled == true {
			for _, col := range value.Collections.ReadAccess {
				if _, ok := c.CollectionResources[col]; !ok {
					value := "CONFIG: One or more API Roots is using a read access collection of " + col + " that is missing from the configuration file"
					c.Logger.Println(value)
					problemsFound++
				}
			}

			for _, col := range value.Collections.WriteAccess {
				if _, ok := c.CollectionResources[col]; !ok {
					value := "CONFIG: One or more API Roots is using a read access collection of " + col + " that is missing from the configuration file"
					c.Logger.Println(value)
					problemsFound++
				}
			}
		}

	} // End for loop on services

	// Log an error if there are no API Root services actually enabled.
	if isServiceEnabled == false {
		c.Logger.Println("CONFIG: While the API Root Server is enabled, there are no API Root Services that are enabled")
		problemsFound++
	}

	// Return errors if there were any
	if problemsFound > 0 {
		c.Logger.Println("ERROR: The API Root configuration has", problemsFound, "error(s)")
	}
	return problemsFound
}
