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

// verifyDiscoveryHTMLConfig - This method will only be called from
// VerifyServerConfig() if DiscoveryServer.HTMLEnabled == true
func (ezt *ServerConfigType) verifyDiscoveryHTMLConfig() error {
	var problemsFound = 0

	if ezt.DiscoveryServer.HTMLTemplateFile == "" {
		log.Println("CONFIG: The Discovery Server is missing the htmltemplatefile directive in the configuration file")
		problemsFound++
	} else {
		// Lets check to make sure the file exists.
		problemsFound += ezt.verifyHTMLFileExists(ezt.DiscoveryServer.HTMLTemplateFile)

		// Need to check to see if the HTML resource file was redefined at each service level
		for i, s := range ezt.DiscoveryServer.Services {
			// Set the HTMLEnabled to true at the service level since it is true
			// at the parent level. We do not allow this to be redefined in the
			// configuration file.
			ezt.DiscoveryServer.Services[i].HTMLEnabled = true

			// Set the HTMLTemplatePath to the prefix + HTMLTemplateDir from the
			// global config. This will make it easier for us to use later on.
			ezt.DiscoveryServer.Services[i].HTMLTemplatePath = ezt.Global.Prefix + "/" + ezt.Global.HTMLTemplateDir

			// If it is not defined at the service level, lets copy from the
			// parent, this will make it easier to work with later on.
			if s.HTMLTemplateFile == "" {
				ezt.DiscoveryServer.Services[i].HTMLTemplateFile = ezt.DiscoveryServer.HTMLTemplateFile
			} else {
				// Only test if the file was redefined at this level. No need to retest the inherited filename since it was already checked
				problemsFound += ezt.verifyHTMLFileExists(s.HTMLTemplateFile)
			}
		} // End for loop
	}

	// Return errors if there were any
	if problemsFound == 1 {
		log.Println("ERROR:", problemsFound, "error was found in the Discovery HTML Template configuration")
		return errors.New("Configuration Errors Found")
	} else if problemsFound > 1 {
		log.Println("ERROR:", problemsFound, "errors were found in the Discovery HTML Template configuration")
		return errors.New("Configuration Errors Found")
	} else {
		return nil
	}
}

// ----------------------------------------
// Verify API Root HTML Files
// ----------------------------------------

// verifyAPIRootHTMLConfig - This method will check each of the defined HTML template files
// and make sure they exist. It will also check to see if any of them have been redefined at
// a service level. If they have, it will check to see if those exists as well.
// retval: problemsFound - an integer representing a count of the number of errors found
func (ezt *ServerConfigType) verifyAPIRootHTMLConfig() error {
	var problemsFound = 0

	if ezt.APIRootServer.HTMLTemplateFiles.APIRoot == "" {
		log.Println("CONFIG: The API Root Service is missing the 'apiroot' directive from the `htmltemplatefiles` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += ezt.verifyHTMLFileExists(ezt.APIRootServer.HTMLTemplateFiles.APIRoot)
	}

	if ezt.APIRootServer.HTMLTemplateFiles.Collections == "" {
		log.Println("CONFIG: The API Root Service is missing the 'collections' directive from the `htmltemplatefiles` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += ezt.verifyHTMLFileExists(ezt.APIRootServer.HTMLTemplateFiles.Collections)
	}

	if ezt.APIRootServer.HTMLTemplateFiles.Collection == "" {
		log.Println("CONFIG: The API Root Service is missing the 'collection' directive from the `htmltemplatefiles` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += ezt.verifyHTMLFileExists(ezt.APIRootServer.HTMLTemplateFiles.Collection)
	}

	if ezt.APIRootServer.HTMLTemplateFiles.Objects == "" {
		log.Println("CONFIG: The API Root Service is missing the 'objects' directive from the `htmltemplatefiles` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += ezt.verifyHTMLFileExists(ezt.APIRootServer.HTMLTemplateFiles.Objects)
	}

	// Lets check to see if any of the HTML template files have been redefined at the service level

	for i, s := range ezt.APIRootServer.Services {

		// Lets set the HTTPEnabled to true at the service level since it
		// is true at the parent level. We do not allow this to be redefined
		// in the configuration file.
		ezt.APIRootServer.Services[i].HTMLEnabled = true

		// Set the HTMLTemplatePath to the prefix + HTMLTemplateDir from the
		// global config. This will make it easier for us to use later on.
		ezt.APIRootServer.Services[i].HTMLTemplatePath = ezt.Global.Prefix + "/" + ezt.Global.HTMLTemplateDir

		// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
		if s.HTMLTemplateFiles.APIRoot == "" {
			ezt.APIRootServer.Services[i].HTMLTemplateFiles.APIRoot = ezt.APIRootServer.HTMLTemplateFiles.APIRoot
		} else {
			// Only test if the file was redefined at ezt level. No need to retest the inherited filename since it was already checked
			problemsFound += ezt.verifyHTMLFileExists(s.HTMLTemplateFiles.APIRoot)
		}

		// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
		if s.HTMLTemplateFiles.Collections == "" {
			ezt.APIRootServer.Services[i].HTMLTemplateFiles.Collections = ezt.APIRootServer.HTMLTemplateFiles.Collections
		} else {
			// Only test if the file was redefined at ezt level. No need to retest the inherited filename since it was already checked
			problemsFound += ezt.verifyHTMLFileExists(s.HTMLTemplateFiles.Collections)
		}

		// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
		if s.HTMLTemplateFiles.Collection == "" {
			ezt.APIRootServer.Services[i].HTMLTemplateFiles.Collection = ezt.APIRootServer.HTMLTemplateFiles.Collection
		} else {
			// Only test if the file was redefined at this level. No need to retest the inherited filename since it was already checked
			problemsFound += ezt.verifyHTMLFileExists(s.HTMLTemplateFiles.Collection)
		}

		// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
		if s.HTMLTemplateFiles.Objects == "" {
			ezt.APIRootServer.Services[i].HTMLTemplateFiles.Objects = ezt.APIRootServer.HTMLTemplateFiles.Objects
		} else {
			// Only test if the file was redefined at this level. No need to retest the inherited filename since it was already checked
			problemsFound += ezt.verifyHTMLFileExists(s.HTMLTemplateFiles.Objects)
		}

	} // End for loop

	// Return errors if there were any
	if problemsFound == 1 {
		log.Println("ERROR:", problemsFound, "error was found in the HTML Template configuration")
		return errors.New("Configuration Errors Found")
	} else if problemsFound > 1 {
		log.Println("ERROR:", problemsFound, "errors were found in the HTML Template configuration")
		return errors.New("Configuration Errors Found")
	} else {
		return nil
	}
}

// -----------------------------------------------------------------------------
// verifyHTMLFileExists - This method will check to make sure the HTML resource file is found on the filesystem
// param: file - a string representing the filename name of the HTML resource file
func (ezt *ServerConfigType) verifyHTMLFileExists(filename string) int {
	filepath := ezt.Global.HTMLTemplateDir + filename
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Println("CONFIG: The HTML template file", filename, "can not be opened:", err)
		return 1
	}
	return 0
}
