// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"errors"
	"log"
	"os"
)

// ----------------------------------------
// Verify Discovery HTML Files
// ----------------------------------------

// verifyDiscoveryHTMLConfig - This method will check each of the defined HTML
// template files and make sure they exist. It will also check to see if any of
// them have been redefined at a service level. If they have, it will check to
// see if those exists as well.
// This method will only be called from VerifyServerConfig() if
// DiscoveryServer.HTMLEnabled == true
func (config *ServerConfigType) verifyDiscoveryHTMLConfig() error {
	var problemsFound = 0

	if config.DiscoveryServer.HTMLBranding.Discovery == "" {
		log.Println("CONFIG: The Discovery Server is missing the htmlbranding.discovery directive in the configuration file")
		problemsFound++
	} else {
		// Lets check to make sure the file exists.
		problemsFound += config.verifyHTMLFileExists(config.DiscoveryServer.HTMLBranding.Discovery)

		// Need to check to see if the HTML resource file was redefined at each service level
		for i, s := range config.DiscoveryServer.Services {
			// Set the HTMLEnabled to true at the service level since it is true
			// at the parent level. We do not allow this to be redefined in the
			// configuration file.
			config.DiscoveryServer.Services[i].HTMLEnabled = true

			// Set the HTMLTemplatePath to the prefix + HTMLTemplateDir from the
			// global config. This will make it easier for us to use later on.
			config.DiscoveryServer.Services[i].HTMLTemplatePath = config.Global.Prefix + config.Global.HTMLTemplateDir

			// If it is not defined at the service level, lets copy from the
			// parent, this will make it easier to work with later on.
			if s.HTMLBranding.Discovery == "" {
				config.DiscoveryServer.Services[i].HTMLBranding.Discovery = config.DiscoveryServer.HTMLBranding.Discovery
			} else {
				// Only test if the file was redefined at this level. No need to
				// retest the inherited filename since it was already checked
				problemsFound += config.verifyHTMLFileExists(s.HTMLBranding.Discovery)
			}
		} // End for loop
	}

	// Return errors if there were any
	if problemsFound > 0 {
		log.Println("ERROR: The Discovery HTML Template configuration has", problemsFound, "error(s)")
		return errors.New("ERROR: Configuration errors found")
	}
	return nil
}

// ----------------------------------------
// Verify API Root HTML Files
// ----------------------------------------

// verifyAPIRootHTMLConfig - This method will check each of the defined HTML template files
// and make sure they exist. It will also check to see if any of them have been redefined at
// a service level. If they have, it will check to see if those exists as well.
// This method will only be called from VerifyServerConfig() if
// APIRootServer.HTMLEnabled == true
func (config *ServerConfigType) verifyAPIRootHTMLConfig() error {
	var problemsFound = 0

	if config.APIRootServer.HTMLBranding.APIRoot == "" {
		log.Println("CONFIG: The API Root Service is missing the 'apiroot' directive from the `htmlbranding.apiroot` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += config.verifyHTMLFileExists(config.APIRootServer.HTMLBranding.APIRoot)
	}

	if config.APIRootServer.HTMLBranding.Collections == "" {
		log.Println("CONFIG: The API Root Service is missing the 'collections' directive from the `htmlbranding.collections` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += config.verifyHTMLFileExists(config.APIRootServer.HTMLBranding.Collections)
	}

	if config.APIRootServer.HTMLBranding.Collection == "" {
		log.Println("CONFIG: The API Root Service is missing the 'collection' directive from the `htmlbranding.collection` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += config.verifyHTMLFileExists(config.APIRootServer.HTMLBranding.Collection)
	}

	if config.APIRootServer.HTMLBranding.Objects == "" {
		log.Println("CONFIG: The API Root Service is missing the 'objects' directive from the `htmlbranding.objects` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += config.verifyHTMLFileExists(config.APIRootServer.HTMLBranding.Objects)
	}

	if config.APIRootServer.HTMLBranding.Manifest == "" {
		log.Println("CONFIG: The API Root Service is missing the 'manifest' directive from the `htmlbranding.objects` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += config.verifyHTMLFileExists(config.APIRootServer.HTMLBranding.Manifest)
	}

	// Lets check to see if any of the HTML template files have been redefined at the service level
	for i, s := range config.APIRootServer.Services {

		// Lets set the HTTPEnabled to true at the service level since it
		// is true at the parent level. We do not allow this to be redefined
		// in the configuration file.
		config.APIRootServer.Services[i].HTMLEnabled = true

		// Set the HTMLTemplatePath to the prefix + HTMLTemplateDir from the
		// global config. This will make it easier for us to use later on.
		config.APIRootServer.Services[i].HTMLTemplatePath = config.Global.Prefix + config.Global.HTMLTemplateDir

		// If it is not defined at the service level, lets copy in the parent,
		// this will make it easier to work with later on
		if s.HTMLBranding.APIRoot == "" {
			config.APIRootServer.Services[i].HTMLBranding.APIRoot = config.APIRootServer.HTMLBranding.APIRoot
		} else {
			// Only test if the file was redefined at config level. No need to
			// retest the inherited filename since it was already checked
			problemsFound += config.verifyHTMLFileExists(s.HTMLBranding.APIRoot)
		}

		// If it is not defined at the service level, lets copy in the parent,
		// this will make it easier to work with later on
		if s.HTMLBranding.Collections == "" {
			config.APIRootServer.Services[i].HTMLBranding.Collections = config.APIRootServer.HTMLBranding.Collections
		} else {
			// Only test if the file was redefined at config level. No need to
			// retest the inherited filename since it was already checked
			problemsFound += config.verifyHTMLFileExists(s.HTMLBranding.Collections)
		}

		// If it is not defined at the service level, lets copy in the parent,
		// this will make it easier to work with later on
		if s.HTMLBranding.Collection == "" {
			config.APIRootServer.Services[i].HTMLBranding.Collection = config.APIRootServer.HTMLBranding.Collection
		} else {
			// Only test if the file was redefined at this level. No need to
			// retest the inherited filename since it was already checked
			problemsFound += config.verifyHTMLFileExists(s.HTMLBranding.Collection)
		}

		// If it is not defined at the service level, lets copy in the parent,
		// this will make it easier to work with later on
		if s.HTMLBranding.Objects == "" {
			config.APIRootServer.Services[i].HTMLBranding.Objects = config.APIRootServer.HTMLBranding.Objects
		} else {
			// Only test if the file was redefined at this level. No need to
			// retest the inherited filename since it was already checked
			problemsFound += config.verifyHTMLFileExists(s.HTMLBranding.Objects)
		}

		// If it is not defined at the service level, lets copy in the parent,
		// this will make it easier to work with later on
		if s.HTMLBranding.Manifest == "" {
			config.APIRootServer.Services[i].HTMLBranding.Manifest = config.APIRootServer.HTMLBranding.Manifest
		} else {
			// Only test if the file was redefined at this level. No need to
			// retest the inherited filename since it was already checked
			problemsFound += config.verifyHTMLFileExists(s.HTMLBranding.Manifest)
		}

	} // End for loop

	// Return errors if there were any
	if problemsFound > 0 {
		log.Println("ERROR: The API Root HTML Template configuration has", problemsFound, "error(s)")
		return errors.New("ERROR: Configuration errors found")
	}
	return nil
}

// -----------------------------------------------------------------------------
// verifyHTMLFileExists - This method will take in a string representing the
// filename of the HTML resource file and check to make sure that HTML resource
// file is found on the filesystem
func (config *ServerConfigType) verifyHTMLFileExists(filename string) int {
	filepath := config.Global.Prefix + config.Global.HTMLTemplateDir + filename
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Println("CONFIG: The HTML template file", filename, "can not be opened:", err)
		return 1
	}
	return 0
}
