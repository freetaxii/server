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

// -----------------------------------------------------------------------------
// verifyConfigDirectives - This method will verify that each required configuration directive is present
// retval: error
func (this *ServerConfigType) verifyConfigDirectives() error {
	var problemsFound int = 0

	problemsFound += this.verifyHtmlTemplateFiles()

	// Listen Directive
	if this.System.Listen == "" {
		log.Println("CONFIG: The listen directive is missing from the configuration file")
		problemsFound++
	}

	// Database Config
	if this.System.DbConfig == true && this.System.DbFile == "" {
		log.Println("CONFIG: The 'dbconfig' directive is set to true, however, the 'dbfile' directive is missing from the configuration file")
		problemsFound++
	}

	// Discovery Service Directive
	for i, _ := range this.DiscoveryService.Services {

		// Verify the Discovery Name is defined in the configuraiton file. This is used as the path name
		if this.DiscoveryService.Services[i].Name == "" {
			log.Println("CONFIG: One or more Discovery Services is missing the 'name' directive in the configuration file")
			problemsFound++
		} else {
			this.DiscoveryService.Services[i].Path = "/" + this.DiscoveryService.Services[i].Name + "/"
		}

		// Verify the Discovery Resource is found
		if _, ok := this.DiscoveryResources[this.DiscoveryService.Services[i].Resource]; !ok {
			value := "CONFIG: The Discovery Resource " + this.DiscoveryService.Services[i].Resource + " is missing from the configuration file"
			log.Println(value)
			problemsFound++
		}
	}

	// API Service Directives
	for i, _ := range this.ApiRootService.Services {

		// Verify the API Name is defined in the configuration file. This is used as the path name
		if this.ApiRootService.Services[i].Name == "" {
			log.Println("CONFIG: One or more API Root Services is missing the 'name' directive in the configuration file")
			problemsFound++
		} else {
			this.ApiRootService.Services[i].Path = "/" + this.ApiRootService.Services[i].Name + "/"
		}

		// Verify the API Resource is found
		if _, ok := this.ApiRootResources[this.ApiRootService.Services[i].Resource]; !ok {
			value := "CONFIG: The API Root Resource " + this.ApiRootService.Services[i].Resource + " is missing from the configuration file"
			log.Println(value)
			problemsFound++
		}
	}

	// Return errors if there were any
	if problemsFound == 1 {
		log.Println("ERROR: ", problemsFound, " error was found in the configuration file")
		return errors.New("Configuration Errors Found")
	} else if problemsFound > 1 {
		log.Println("ERROR: ", problemsFound, " errors were found in the configuration file")
		return errors.New("Configuration Errors Found")
	} else {
		return nil
	}
}

// -----------------------------------------------------------------------------
// verifyHtmlTemplateFiles - This method will check each of the defined HTML template files
// and make sure they exist. It will also check to see if any of them have been redefined at
// a service level. If they have, it will check to see if those exists as well.
// retval: problemsFound - an integer representing a count of the number of errors found
func (this *ServerConfigType) verifyHtmlTemplateFiles() int {
	var problemsFound int = 0

	// Discovery HTML File
	if this.DiscoveryService.Html == true && this.DiscoveryService.HtmlFile == "" {
		log.Println("CONFIG: The Discovery Service is missing the 'htmlfile' directive in the configuration file")
		problemsFound++
	} else if this.DiscoveryService.Html == true && this.DiscoveryService.HtmlFile != "" {

		// Lets check to make sure the file exists.
		problemsFound += this.verifyHtmlFileExists(this.DiscoveryService.HtmlFile)

		// Need to check to see if the HTML resource file was redefined at each service level
		for i, _ := range this.DiscoveryService.Services {
			// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
			if this.DiscoveryService.Services[i].HtmlFile == "" {
				this.DiscoveryService.Services[i].HtmlFile = this.DiscoveryService.HtmlFile
			} else {
				// Only test if the file was redefined at this level. No need to retest the inhertied filename since it was already checked
				problemsFound += this.verifyHtmlFileExists(this.DiscoveryService.Services[i].HtmlFile)
			}
		} // End for loop
	} // End If Discovery HTML File

	// ----------------------------------------
	// Verify API Root HTML Files
	// ----------------------------------------
	if this.ApiRootService.Html == true && this.ApiRootService.HtmlFiles.ApiRoot == "" {
		log.Println("CONFIG: The API Root Service is missing the 'apiroot' directive from the `htmlfiles` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += this.verifyHtmlFileExists(this.ApiRootService.HtmlFiles.ApiRoot)
	}

	if this.ApiRootService.Html == true && this.ApiRootService.HtmlFiles.Collections == "" {
		log.Println("CONFIG: The API Root Service is missing the 'collections' directive from the `htmlfiles` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += this.verifyHtmlFileExists(this.ApiRootService.HtmlFiles.Collections)
	}

	if this.ApiRootService.Html == true && this.ApiRootService.HtmlFiles.Collection == "" {
		log.Println("CONFIG: The API Root Service is missing the 'collection' directive from the `htmlfiles` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += this.verifyHtmlFileExists(this.ApiRootService.HtmlFiles.Collection)
	}

	if this.ApiRootService.Html == true && this.ApiRootService.HtmlFiles.Objects == "" {
		log.Println("CONFIG: The API Root Service is missing the 'objects' directive from the `htmlfiles` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += this.verifyHtmlFileExists(this.ApiRootService.HtmlFiles.Objects)
	}

	// Lets check to see if any of the HTML template files have been redefined at the service level
	if this.ApiRootService.Html == true {

		for i, _ := range this.ApiRootService.Services {

			// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
			if this.ApiRootService.Services[i].HtmlFiles.ApiRoot == "" {
				this.ApiRootService.Services[i].HtmlFiles.ApiRoot = this.ApiRootService.HtmlFiles.ApiRoot
			} else {
				// Only test if the file was redefined at this level. No need to retest the inhertied filename since it was already checked
				problemsFound += this.verifyHtmlFileExists(this.ApiRootService.Services[i].HtmlFiles.ApiRoot)
			}

			// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
			if this.ApiRootService.Services[i].HtmlFiles.Collections == "" {
				this.ApiRootService.Services[i].HtmlFiles.Collections = this.ApiRootService.HtmlFiles.Collections
			} else {
				// Only test if the file was redefined at this level. No need to retest the inhertied filename since it was already checked
				problemsFound += this.verifyHtmlFileExists(this.ApiRootService.Services[i].HtmlFiles.Collections)
			}

			// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
			if this.ApiRootService.Services[i].HtmlFiles.Collection == "" {
				this.ApiRootService.Services[i].HtmlFiles.Collection = this.ApiRootService.HtmlFiles.Collection
			} else {
				// Only test if the file was redefined at this level. No need to retest the inhertied filename since it was already checked
				problemsFound += this.verifyHtmlFileExists(this.ApiRootService.Services[i].HtmlFiles.Collection)
			}

			// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
			if this.ApiRootService.Services[i].HtmlFiles.Objects == "" {
				this.ApiRootService.Services[i].HtmlFiles.Objects = this.ApiRootService.HtmlFiles.Objects
			} else {
				// Only test if the file was redefined at this level. No need to retest the inhertied filename since it was already checked
				problemsFound += this.verifyHtmlFileExists(this.ApiRootService.Services[i].HtmlFiles.Objects)
			}

		} // End for loop
	} // End If API Root HTML File
	return problemsFound
}

// -----------------------------------------------------------------------------
// verifyHtmlFileExists - This method will check to make sure the HTML resource file is found on the filesystem
// param: file - a string representing the filename name of the HTML resource file
func (this *ServerConfigType) verifyHtmlFileExists(filename string) int {
	filepath := this.System.HtmlDir + filename
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Println("CONFIG: The HTML template file can not be opened", err)
		return 1
	}
	return 0
}
