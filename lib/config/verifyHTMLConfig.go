// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package config

import (
	"strconv"
	"strings"

	"github.com/gologme/log"
)

// ----------------------------------------
// Verify Discovery HTML Files
// ----------------------------------------

/*
verifyGlobalHTMLConfig - This method will check each of the defined HTML settings
in the global configuration and return the number of errors found.
*/
func (c *ServerConfigType) verifyGlobalHTMLConfig() int {
	var problemsFound = 0

	// ----------------------------------------------------------------------
	// Verify TemplateDir is defined and exists on the file system
	// ----------------------------------------------------------------------
	// Check to see if the template directory is implicitly set to null ("") or
	// explicitly set to null ("null"). When set explicitly to null that means
	// it is also invalid.
	if c.HTML.TemplateDir.Value == "" || c.HTML.TemplateDir.Valid == false {
		log.Infoln("CONFIG: The global HTML configuration is missing the html.templatedir directive in the configuration file")
		problemsFound++
	} else {
		problemsFound += c.verifyHTMLTemplateDir("html.templatedir", c.HTML.TemplateDir)
		c.HTML.FullTemplatePath = c.Global.Prefix + c.HTML.TemplateDir.Value
	}

	// ----------------------------------------------------------------------
	// Verify actual template files are defined and exist on the file system
	// ----------------------------------------------------------------------
	problemsFound += c.verifyGlobalHTMLTemplateFile("html.templatefiles.discovery", c.HTML.FullTemplatePath, c.HTML.TemplateFiles.Discovery)
	problemsFound += c.verifyGlobalHTMLTemplateFile("html.templatefiles.apiroot", c.HTML.FullTemplatePath, c.HTML.TemplateFiles.APIRoot)
	problemsFound += c.verifyGlobalHTMLTemplateFile("html.templatefiles.collections", c.HTML.FullTemplatePath, c.HTML.TemplateFiles.Collections)
	problemsFound += c.verifyGlobalHTMLTemplateFile("html.templatefiles.collection", c.HTML.FullTemplatePath, c.HTML.TemplateFiles.Collection)
	problemsFound += c.verifyGlobalHTMLTemplateFile("html.templatefiles.objects", c.HTML.FullTemplatePath, c.HTML.TemplateFiles.Objects)

	// ----------------------------------------------------------------------
	// Return number of errors if there are any
	// ----------------------------------------------------------------------
	if problemsFound > 0 {
		log.Println("ERROR: The global HTML configuration has", problemsFound, "error(s)")
	}
	return problemsFound
}

/*
verifyHTMLTemplateDir - This method will verify that the template directory in
the configuration has a trailing slash and that it is found on the file system.
*/
func (c *ServerConfigType) verifyHTMLTemplateDir(configPath string, templateDir JSONstring) int {
	var problemsFound = 0

	if !strings.HasSuffix(templateDir.Value, "/") {
		log.Println("CONFIG: The" + configPath + "directive is missing the ending slash '/'")
		problemsFound++
	}

	filepath := c.Global.Prefix + templateDir.Value
	if !c.exists(filepath) {
		log.Infoln("CONFIG: The HTML template path", filepath, "can not be opened")
		problemsFound++
	}
	return problemsFound
}

/*
verifyGlobalHTMLTemplateFile - This method will verify that the global HTML
template files are defined and that they can be found on the file system.
*/
func (c *ServerConfigType) verifyGlobalHTMLTemplateFile(configPath, templatePath string, template JSONstring) int {
	var problemsFound = 0

	if templatePath == "" {
		log.Infoln("CONFIG: The HTML template path used by" + configPath + "is missing")
		problemsFound++
		return problemsFound
	}

	if template.Value == "" || template.Valid == false {
		log.Infoln("CONFIG: The HTML configuration is missing the" + configPath + "directive in the configuration file")
		problemsFound++
		return problemsFound
	}

	problemsFound += c.verifyHTMLTemplateFile(configPath, templatePath, template)
	return problemsFound
}

/*
verifyHTMLTemplateFile - This method will verify that HTML template files are
found on the file system.
*/
func (c *ServerConfigType) verifyHTMLTemplateFile(configPath, templatePath string, template JSONstring) int {
	var problemsFound = 0

	if templatePath == "" {
		log.Infoln("CONFIG: The HTML template path used by" + configPath + "is missing")
		problemsFound++
		return problemsFound
	}

	filepath := templatePath + template.Value
	if !c.exists(filepath) {
		log.Infoln("CONFIG: The HTML template path", filepath, "defined at"+configPath+"can not be opened")
		problemsFound++
	}

	return problemsFound
}

/*
verifyDiscoveryHTMLConfig - This method will check each of the services to see
if the HTML configuration has been redefined. If the values have not been
redefined then the global settings will be copied in to this level. The actual
HTTP handlers will use the settings found in these services and not the global
settings.

This method is called from serverConfig.go-verifyServerConfig()
*/
func (c *ServerConfigType) verifyDiscoveryHTMLConfig() int {
	var problemsFound = 0

	// If HTML output is not enabled globally, then skip these tests
	if c.HTML.Enabled.Value == false {
		return problemsFound
	}

	// Check to see if any of the HTML configurations were redefined at each service level
	for i, s := range c.DiscoveryServer.Services {
		indexString := strconv.Itoa(i)

		// Check to see if the following values were redefined and valid. If
		// they were not redefined or they are invalid (set to "null" and thus
		// invalid) then lets set them to the same as the global configuration.
		// Copy all of the settings for the object, not just the value.
		// If the value was redefined and is valid, lets just leave it alone.
		if s.HTML.Enabled.Set == false || s.HTML.Enabled.Valid == false {
			c.DiscoveryServer.Services[i].HTML.Enabled = c.HTML.Enabled
		}

		if s.HTML.TemplateDir.Set == false || s.HTML.TemplateDir.Valid == false {
			c.DiscoveryServer.Services[i].HTML.TemplateDir = c.HTML.TemplateDir
			c.DiscoveryServer.Services[i].HTML.FullTemplatePath = c.HTML.FullTemplatePath
		} else {
			// If it was redefined we need to update the TemplatePath, first lets
			// verify that the template directory is found on the file system,
			// if it is then copy all of the settings for the template path and
			// then update the actual value.
			text := "discoveryserver.services[" + indexString + "].html.templatedir"
			problemsFound += c.verifyHTMLTemplateDir(text, s.HTML.TemplateDir)
			c.DiscoveryServer.Services[i].HTML.FullTemplatePath = c.HTML.FullTemplatePath
			c.DiscoveryServer.Services[i].HTML.FullTemplatePath = c.Global.Prefix + s.HTML.TemplateDir.Value
		}

		if s.HTML.TemplateFiles.Discovery.Set == false || s.HTML.TemplateFiles.Discovery.Valid == false {
			c.DiscoveryServer.Services[i].HTML.TemplateFiles.Discovery = c.HTML.TemplateFiles.Discovery
		} else {
			// If it was redefined we need to verify that it is found on the file system.
			text := "discoveryserver.services[" + indexString + "].html.templatefiles.discovery"
			problemsFound += c.verifyHTMLTemplateFile(text, s.HTML.FullTemplatePath, s.HTML.TemplateFiles.Discovery)
		}
	} // End for loop

	// ----------------------------------------------------------------------
	// Return number of errors if there are any
	// ----------------------------------------------------------------------
	if problemsFound > 0 {
		log.Println("ERROR: The Discovery HTML configuration has", problemsFound, "error(s)")
	}
	return problemsFound
}

// ----------------------------------------
// Verify API Root HTML Files
// ----------------------------------------

// verifyAPIRootHTMLConfig - This method will check each of the defined HTML template files
// and make sure they exist. It will also check to see if any of them have been redefined at
// a service level. If they have, it will check to see if those exists as well.
// This method will only be called from VerifyServerConfig() if
// APIRootServer.HTMLEnabled == true
func (c *ServerConfigType) verifyAPIRootHTMLConfig() int {
	var problemsFound = 0

	// If HTML output is not enabled globally, then skip these tests
	if c.HTML.Enabled.Value == false {
		return problemsFound
	}

	// Check to see if any of the HTML configurations were redefined at each service level
	for i, s := range c.APIRootServer.Services {
		indexString := strconv.Itoa(i)

		// Check to see if the following values were redefined and valid. If
		// they were not redefined or they are invalid (set to "null" and thus
		// invalid) then lets set them to the same as the global configuration.
		// Copy all of the settings for the object, not just the value.
		// If the value was redefined and is valid, lets just leave it alone.
		if s.HTML.Enabled.Set == false || s.HTML.Enabled.Valid == false {
			c.APIRootServer.Services[i].HTML.Enabled = c.HTML.Enabled
		}

		if s.HTML.TemplateDir.Set == false || s.HTML.TemplateDir.Valid == false {
			c.APIRootServer.Services[i].HTML.TemplateDir = c.HTML.TemplateDir
			c.APIRootServer.Services[i].HTML.FullTemplatePath = c.HTML.FullTemplatePath
		} else {
			// If it was redefined we need to update the TemplatePath, first lets
			// verify that the template directory is found on the file system,
			// if it is then copy all of the settings for the template path and
			// then update the actual value.
			text := "apirootserver.services[" + indexString + "].html.templatedir"
			problemsFound += c.verifyHTMLTemplateDir(text, s.HTML.TemplateDir)
			c.APIRootServer.Services[i].HTML.FullTemplatePath = c.HTML.FullTemplatePath
			c.APIRootServer.Services[i].HTML.FullTemplatePath = c.Global.Prefix + s.HTML.TemplateDir.Value
		}

		if s.HTML.TemplateFiles.APIRoot.Set == false || s.HTML.TemplateFiles.APIRoot.Valid == false {
			c.APIRootServer.Services[i].HTML.TemplateFiles.APIRoot = c.HTML.TemplateFiles.APIRoot
		} else {
			// If it was redefined we need to verify that it is found on the file system.
			text := "apirootserver.services[" + indexString + "].html.templatefiles.discovery"
			problemsFound += c.verifyHTMLTemplateFile(text, s.HTML.FullTemplatePath, s.HTML.TemplateFiles.APIRoot)
		}

		if s.HTML.TemplateFiles.Collections.Set == false || s.HTML.TemplateFiles.Collections.Valid == false {
			c.APIRootServer.Services[i].HTML.TemplateFiles.Collections = c.HTML.TemplateFiles.Collections
		} else {
			// If it was redefined we need to verify that it is found on the file system.
			text := "apirootserver.services[" + indexString + "].html.templatefiles.discovery"
			problemsFound += c.verifyHTMLTemplateFile(text, s.HTML.FullTemplatePath, s.HTML.TemplateFiles.Collections)
		}

		if s.HTML.TemplateFiles.Collection.Set == false || s.HTML.TemplateFiles.Collection.Valid == false {
			c.APIRootServer.Services[i].HTML.TemplateFiles.Collection = c.HTML.TemplateFiles.Collection
		} else {
			// If it was redefined we need to verify that it is found on the file system.
			text := "apirootserver.services[" + indexString + "].html.templatefiles.discovery"
			problemsFound += c.verifyHTMLTemplateFile(text, s.HTML.FullTemplatePath, s.HTML.TemplateFiles.Collection)
		}

		if s.HTML.TemplateFiles.Objects.Set == false || s.HTML.TemplateFiles.Objects.Valid == false {
			c.APIRootServer.Services[i].HTML.TemplateFiles.Objects = c.HTML.TemplateFiles.Objects
		} else {
			// If it was redefined we need to verify that it is found on the file system.
			text := "apirootserver.services[" + indexString + "].html.templatefiles.discovery"
			problemsFound += c.verifyHTMLTemplateFile(text, s.HTML.FullTemplatePath, s.HTML.TemplateFiles.Objects)
		}

		if s.HTML.TemplateFiles.Manifest.Set == false || s.HTML.TemplateFiles.Manifest.Valid == false {
			c.APIRootServer.Services[i].HTML.TemplateFiles.Manifest = c.HTML.TemplateFiles.Manifest
		} else {
			// If it was redefined we need to verify that it is found on the file system.
			text := "apirootserver.services[" + indexString + "].html.templatefiles.discovery"
			problemsFound += c.verifyHTMLTemplateFile(text, s.HTML.FullTemplatePath, s.HTML.TemplateFiles.Manifest)
		}
	} // End for loop

	// ----------------------------------------------------------------------
	// Return number of errors if there are any
	// ----------------------------------------------------------------------
	if problemsFound > 0 {
		log.Println("ERROR: The API Root HTML configuration has", problemsFound, "error(s)")
	}
	return problemsFound
}
