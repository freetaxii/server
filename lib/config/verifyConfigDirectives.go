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
	"strconv"
)

// verifyConfigDirectives - This method will verify that each required configuration directive is present
// retval: error
func (this *ServerConfigType) verifyConfigDirectives() error {
	var errorsFound int = 0

	// Listen Directive
	if this.System.Listen == "" {
		log.Println("CONFIG: The listen directive is missing from the configuration file")
		errorsFound++
	}

	// Database Config
	if this.System.DbConfig == true && this.System.DbFile == "" {
		log.Println("CONFIG: The 'dbconfig' directive is set to true, however, the 'dbfile' directive is missing from the configuration file")
		errorsFound++
	}

	// Discovery HTML File
	if this.DiscoveryService.Html == true && this.DiscoveryService.HtmlFile == "" {
		log.Println("CONFIG: The Discovery Service is missing the 'htmlfile' directive in the configuration file")
		errorsFound++
	} else if this.DiscoveryService.Html == true && this.DiscoveryService.HtmlFile != "" {

		// Lets check to make sure the file exists.
		filename := this.System.HtmlDir + this.DiscoveryService.HtmlFile
		err := this.verifyHtmlFileExists(filename)
		if err != nil {
			log.Println("CONFIG: The Discovery Service HTML template file can not be opened", err)
			errorsFound++
		}

		// Need to check to see if the HTML resource file was redefined at each service level
		for i, _ := range this.DiscoveryService.Services {
			// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
			if this.DiscoveryService.Services[i].HtmlFile == "" {
				this.DiscoveryService.Services[i].HtmlFile = this.DiscoveryService.HtmlFile
			} else {
				// Only test if the file was redefined at this level. No need to retest the inhertied filename since it was already checked
				filename := this.System.HtmlDir + this.DiscoveryService.Services[i].HtmlFile
				err := this.verifyHtmlFileExists(filename)
				if err != nil {
					log.Println("CONFIG: The Discovery Service HTML template file can not be opened", err)
					errorsFound++
				}
			}
		} // End for loop
	} // End If Discovery HTML File

	// API Root HTML File
	if this.ApiRootService.Html == true && this.ApiRootService.HtmlFile == "" {
		log.Println("CONFIG: The API Root Service is missing the 'htmlfile' directive in the configuration file")
		errorsFound++
	} else if this.ApiRootService.Html == true && this.ApiRootService.HtmlFile != "" {

		// Lets check to make sure the file exists
		filename := this.System.HtmlDir + this.ApiRootService.HtmlFile
		err := this.verifyHtmlFileExists(filename)
		if err != nil {
			log.Println("CONFIG: The API Root Service HTML template file can not be opened", err)
			errorsFound++
		}

		// Need to check to see if the HTML resource file was redefined at each service level
		for i, _ := range this.ApiRootService.Services {
			var err error

			// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
			if this.ApiRootService.Services[i].HtmlFile == "" {
				this.ApiRootService.Services[i].HtmlFile = this.ApiRootService.HtmlFile
			} else {
				// Only test if the file was redefined at this level. No need to retest the inhertied filename since it was already checked
				filename := this.System.HtmlDir + this.ApiRootService.Services[i].HtmlFile
				err = this.verifyHtmlFileExists(filename)
				if err != nil {
					log.Println("CONFIG: The API Root Service HTML template file can not be opened", err)
					errorsFound++
				}
			}

			// Only test if the file was redefined at this level. No need to reset the inherited filename since it was already checked
			collectionsHtmlfilename := this.System.HtmlDir + this.ApiRootService.Services[i].Collections.HtmlFile
			err = this.verifyHtmlFileExists(collectionsHtmlfilename)
			if err != nil {
				log.Println("CONFIG: The Collections Service HTML template file can not be opened", err)
				errorsFound++
			}

			collectionHtmlfilename := this.System.HtmlDir + this.ApiRootService.Services[i].Collection.HtmlFile
			err = this.verifyHtmlFileExists(collectionHtmlfilename)
			if err != nil {
				log.Println("CONFIG: The Collection Service HTML template file can not be opened", err)
				errorsFound++
			}

			objectsHtmlfilename := this.System.HtmlDir + this.ApiRootService.Services[i].Objects.HtmlFile
			err = this.verifyHtmlFileExists(objectsHtmlfilename)
			if err != nil {
				log.Println("CONFIG: The Objects Service HTML template file can not be opened", err)
				errorsFound++
			}
		} // End for loop
	} // End If API Root HTML File

	// Discovery Service Directive
	for i, _ := range this.DiscoveryService.Services {

		// Verify the Discovery Name is defined in the configuraiton file. This is used as the path name
		if this.DiscoveryService.Services[i].Name == "" {
			log.Println("CONFIG: One or more Discovery Services is missing the 'name' directive in the configuration file")
			errorsFound++
		} else {
			this.DiscoveryService.Services[i].Path = "/" + this.DiscoveryService.Services[i].Name + "/"
		}

		// Verify the Discovery Resource is found
		if _, ok := this.DiscoveryResources[this.DiscoveryService.Services[i].Resource]; !ok {
			value := "CONFIG: The Discovery Resource " + this.DiscoveryService.Services[i].Resource + " is missing from the configuration file"
			log.Println(value)
			errorsFound++
		}
	}

	// API Service Directives
	for i, _ := range this.ApiRootService.Services {

		// Verify the API Name is defined in the configuration file. This is used as the path name
		if this.ApiRootService.Services[i].Name == "" {
			log.Println("CONFIG: One or more API Root Services is missing the 'name' directive in the configuration file")
			errorsFound++
		} else {
			this.ApiRootService.Services[i].Path = "/" + this.ApiRootService.Services[i].Name + "/"
		}

		// Verify the API Resource is found
		if _, ok := this.ApiRootResources[this.ApiRootService.Services[i].Resource]; !ok {
			value := "CONFIG: The API Root Resource " + this.ApiRootService.Services[i].Resource + " is missing from the configuration file"
			log.Println(value)
			errorsFound++
		}

		if this.ApiRootService.Services[i].Collections.HtmlFile == "" {
			log.Println("CONFIG: The Collections Service is missing the 'htmlfile' directive in the configuration file")
			errorsFound++
		}
	}

	// Return errors if there were any
	if errorsFound == 1 {
		value := "ERROR: " + strconv.Itoa(errorsFound) + " error was found in the configuration file"
		return errors.New(value)
	} else if errorsFound > 1 {
		value := "ERROR: " + strconv.Itoa(errorsFound) + " errors were found in the configuration file"
		return errors.New(value)
	} else {
		return nil
	}
}

// verifyHtmlFileExists - This method will check to make sure the HTML resource file is found on the filesystem
// param: file - a string representing the full / relative filename path to the HTML resource file
func (this *ServerConfigType) verifyHtmlFileExists(file string) error {
	filename := this.System.HtmlDir + file
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return errors.New("The following HTML template file is not found on the filesystem:" + filename)
	}
	return nil
}
