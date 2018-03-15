// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package config

import (
	"errors"
	"strings"

	"github.com/gologme/log"
)

// ----------------------------------------
// Verify Discovery HTML Files
// ----------------------------------------

/*
verifyHTMLConfig - This method will verify the configuration settings for the
system level HTML configuration options.
*/
func (config *ServerConfigType) verifyHTMLConfig() error {

	// If HTML output is turned off globally, then there no need to check the
	// configuration and verify everything is present and valid.
	if config.HTML.Enabled == false {
		log.Infoln("CONFIG: The global configuration is not configured to use HTML output")
		return nil
	}

	var problemsFound = 0

	// ----------------------------------------------------------------------
	//
	// Verify TemplateDir is defined and exists on the file system
	// If everything is okay, then lets assign the full path to the
	// TemplatePath directive.
	//
	// ----------------------------------------------------------------------

	if config.HTML.TemplateDir == "" {
		log.Infoln("CONFIG: The HTML configuration is missing from the html.templatedir directive in the configuration file")
		problemsFound++
	} else {
		filepath := config.Global.Prefix + config.HTML.TemplateDir

		if !strings.HasSuffix(config.HTML.TemplateDir, "/") {
			log.Println("CONFIG: The html.templatedir directive is missing the ending slash '/'")
			problemsFound++
		}

		if !config.exists(filepath) {
			log.Infoln("CONFIG: The HTML template path", filepath, "can not be opened")
			problemsFound++
		} else {
			// Prefix + Template Directory
			config.HTML.TemplatePath = filepath
		}
	}

	// ----------------------------------------------------------------------
	//
	// Verify actual template files are defined and exist on the file system
	//
	// ----------------------------------------------------------------------

	if config.HTML.TemplateFiles.Discovery == "" {
		log.Infoln("CONFIG: The HTML configuration is missing the html.templatefiles.discovery directive in the configuration file")
		problemsFound++
	} else {
		// Lets check to make sure the file exists.
		problemsFound += config.verifyHTMLTemplateFileExists(config.HTML.TemplateFiles.Discovery)
	}

	if config.HTML.TemplateFiles.APIRoot == "" {
		log.Infoln("CONFIG: The HTML configuration is missing the html.templatefiles.apiroot directive in the configuration file")
		problemsFound++
	} else {
		// Lets check to make sure the file exists.
		problemsFound += config.verifyHTMLTemplateFileExists(config.HTML.TemplateFiles.APIRoot)
	}

	if config.HTML.TemplateFiles.Collections == "" {
		log.Infoln("CONFIG: The HTML configuration is missing the html.templatefiles.collections directive in the configuration file")
		problemsFound++
	} else {
		// Lets check to make sure the file exists.
		problemsFound += config.verifyHTMLTemplateFileExists(config.HTML.TemplateFiles.Collections)
	}

	if config.HTML.TemplateFiles.Collection == "" {
		log.Infoln("CONFIG: The HTML configuration is missing the html.templatefiles.collection directive in the configuration file")
		problemsFound++
	} else {
		// Lets check to make sure the file exists.
		problemsFound += config.verifyHTMLTemplateFileExists(config.HTML.TemplateFiles.Collection)
	}

	if config.HTML.TemplateFiles.Objects == "" {
		log.Infoln("CONFIG: The HTML configuration is missing the html.templatefiles.objects directive in the configuration file")
		problemsFound++
	} else {
		// Lets check to make sure the file exists.
		problemsFound += config.verifyHTMLTemplateFileExists(config.HTML.TemplateFiles.Objects)
	}

	// ----------------------------------------------------------------------
	//
	// Return errors if there are any
	//
	// ----------------------------------------------------------------------

	if problemsFound > 0 {
		log.Println("ERROR: The system level HTML configuration has", problemsFound, "error(s)")
		return errors.New("ERROR: Configuration errors found")
	}
}

// verifyDiscoveryHTMLConfig - This method will check each of the defined HTML
// template files and make sure they exist. It will also check to see if any of
// them have been redefined at a service level. If they have, it will check to
// see if those exists as well.
// This method will only be called from VerifyServerConfig() if
// DiscoveryServer.HTMLEnabled == true
func (config *ServerConfigType) verifyDiscoveryHTMLConfig() error {

	// If HTML output is not enabled globally, then skip these tests
	if config.HTML.Enabled == false {
		return nil
	}

	var problemsFound = 0

	// Check to see if any of the HTML configurations were redefined at each service level
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
			problemsFound += config.verifyHTMLTemplateFileExists(s.HTMLBranding.Discovery)
		}
	} // End for loop

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

	if config.APIRootServer.HTMLBranding.Manifest == "" {
		log.Infoln("CONFIG: The API Root Service is missing the 'manifest' directive from the `htmlbranding.objects` directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += config.verifyHTMLTemplateFileExists(config.APIRootServer.HTMLBranding.Manifest)
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
			problemsFound += config.verifyHTMLTemplateFileExists(s.HTMLBranding.APIRoot)
		}

		// If it is not defined at the service level, lets copy in the parent,
		// this will make it easier to work with later on
		if s.HTMLBranding.Collections == "" {
			config.APIRootServer.Services[i].HTMLBranding.Collections = config.APIRootServer.HTMLBranding.Collections
		} else {
			// Only test if the file was redefined at config level. No need to
			// retest the inherited filename since it was already checked
			problemsFound += config.verifyHTMLTemplateFileExists(s.HTMLBranding.Collections)
		}

		// If it is not defined at the service level, lets copy in the parent,
		// this will make it easier to work with later on
		if s.HTMLBranding.Collection == "" {
			config.APIRootServer.Services[i].HTMLBranding.Collection = config.APIRootServer.HTMLBranding.Collection
		} else {
			// Only test if the file was redefined at this level. No need to
			// retest the inherited filename since it was already checked
			problemsFound += config.verifyHTMLTemplateFileExists(s.HTMLBranding.Collection)
		}

		// If it is not defined at the service level, lets copy in the parent,
		// this will make it easier to work with later on
		if s.HTMLBranding.Objects == "" {
			config.APIRootServer.Services[i].HTMLBranding.Objects = config.APIRootServer.HTMLBranding.Objects
		} else {
			// Only test if the file was redefined at this level. No need to
			// retest the inherited filename since it was already checked
			problemsFound += config.verifyHTMLTemplateFileExists(s.HTMLBranding.Objects)
		}

		// If it is not defined at the service level, lets copy in the parent,
		// this will make it easier to work with later on
		if s.HTMLBranding.Manifest == "" {
			config.APIRootServer.Services[i].HTMLBranding.Manifest = config.APIRootServer.HTMLBranding.Manifest
		} else {
			// Only test if the file was redefined at this level. No need to
			// retest the inherited filename since it was already checked
			problemsFound += config.verifyHTMLTemplateFileExists(s.HTMLBranding.Manifest)
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
// verifyHTMLTemplateFileExists - This method will take in a string representing the
// filename of the HTML resource file and check to make sure that HTML resource
// file is found on the filesystem
func (config *ServerConfigType) verifyHTMLTemplateFileExists(filename string) int {
	filepath := config.Global.Prefix + config.Global.HTMLTemplateDir + filename

	if !config.exists(filepath) {
		log.Infoln("CONFIG: The HTML template file", filename, "can not be opened")
		return 1
	}
	return 0
}
