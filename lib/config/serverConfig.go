// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/freetaxii/libstix2/resources"
	"github.com/gologme/log"
	"github.com/gorilla/mux"
)

/*
ServerConfigType - This type defines the configuration for the entire server.
*/
type ServerConfigType struct {
	Router *mux.Router
	Logger *log.Logger
	Global struct {
		Prefix             string
		Listen             string
		Protocol           string
		TLSDir             string
		TLSKey             string
		TLSCrt             string
		DbConfig           bool
		DbType             string
		DbFile             string
		MaxNumberOfObjects int
	}
	HTML struct {
		HTMLConfigType
	}
	Logging struct {
		Enabled bool
		LogFile string
	}
	DiscoveryServer struct {
		Enabled  bool
		Services []DiscoveryServiceType
	}
	APIRootServer struct {
		Enabled  bool
		Services []APIRootServiceType
	}
	DiscoveryResources  map[string]resources.DiscoveryType  // The key in the map is the ResourceID
	APIRootResources    map[string]resources.APIRootType    // The key in the map is the ResourceID
	CollectionResources map[string]resources.CollectionType // The key in the map is the ResourceID
}

/*
BaseServiceType - This struct represents the common properties between the
Discovery and API-Root services.

Path          - The URL path for this service
Enabled       - Is this service enabled
ResourceID    - A unique ID for the resource that this service is using
ResourcePath  - The actual full URL path for the resource, used for the handler to know where to listen.
HTML          - The configuration for generating HTML output
*/
type BaseServiceType struct {
	Enabled    bool           // User defined in configuration file
	Path       string         // User defined in configuration file
	ResourceID string         // User defined in configuration file
	HTML       HTMLConfigType // User defined in configuration file or set in the verify scripts.
	FullPath   string         // Set in verifyDiscoveryConfig() or verifyAPIRootConfig()
}

/*
DiscoveryServiceType - This struct represents an instance of a Discovery server.
If someone tries to set the 'resourcepath' directive in the configuration file it
will get overwritten in code.
*/
type DiscoveryServiceType struct {
	BaseServiceType
}

/*
APIRootServiceType - This struct represents an instance of an API Root server.
If someone tries to set the 'path' directive in the configuration file it
will just get overwritten in code.
*/
type APIRootServiceType struct {
	BaseServiceType
	Collections struct {
		Enabled     bool     // User defined in configuration file
		ResourceIDs []string // User defined in configuration file. A list of collections that are members of this API Root
		FullPath    string   // Set in verifyAPIRootConfig()
	}
}

/*
HTMLConfigType - This struct holds the configuration elements for generating HTML
output. This is used at the top level of the configuration file as well as in
each individual service. This means individual services can have a different
HTML configuration. I needed to setup my own types for JSON boolean and strings
since leaving it blank at a child level, would have equaled "false" or "". This
would have been equivalent to turning it off, which is not what is wanted. Leaving
it blank would mean to inherit from the parent. But since Go is a strictly typed
language, you need to create a type that can handle that case.

Enabled       - Is HTML enabled for this service
TemplateDir   - The location of the template files relative to the base of the application (prefix)
TemplatePath  - The full path of the template directory (prefix + TemplateDir)
TemplateFiles - The HTML template filenames in the template directory for the following services
*/
type HTMLConfigType struct {
	Enabled       JSONbool   // User defined in configuration file or set in verifyHTMLConfig()
	TemplateDir   JSONstring // User defined in configuration file or set in verifyHTMLConfig()
	TemplateFiles struct {
		Discovery   JSONstring // User defined in configuration file or set in verifyHTMLConfig()
		APIRoot     JSONstring
		Collections JSONstring
		Collection  JSONstring
		Objects     JSONstring
		Manifest    JSONstring
	}
	FullTemplatePath string // Set in verifyHTMLConfig(), this is the full path to template files
}

// ----------------------------------------------------------------------
//
// Public Create Functions
//
// ----------------------------------------------------------------------

/*
New - This function will return a ServerConfigType, load the current configuration
from a file, and verify that the configuration is correct.
*/
func New(logger *log.Logger, filename string) (ServerConfigType, error) {
	var c ServerConfigType

	if logger == nil {
		c.Logger = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		c.Logger = logger
	}

	err1 := c.loadServerConfig(filename)
	if err1 != nil {
		return c, err1
	}

	// In addition to checking the configuration for completeness the verify
	// process will also populate some of the helper values.
	err2 := c.verifyServerConfig()
	if err2 != nil {
		return c, err2
	}
	return c, nil
}

// --------------------------------------------------
//
// Load Configuration File, Parse JSON, and Verify
//
// --------------------------------------------------

/*
loadServerConfig - This methods takes in a string value representing a
filename of the configuration file and loads the configuration into memory.
*/
func (c *ServerConfigType) loadServerConfig(filename string) error {
	// TODO - Need to make make a validation check for the configuration file

	// Open and read configuration file
	sysConfigFileData, err1 := os.Open(filename)
	defer sysConfigFileData.Close()
	if err1 != nil {
		return fmt.Errorf("error opening configuration file: %v", err1)
	}

	// --------------------------------------------------
	// Decode JSON configuration file
	// --------------------------------------------------
	// Use decoder instead of unmarshal so we can handle stream data
	decoder := json.NewDecoder(sysConfigFileData)
	err2 := decoder.Decode(c)

	if err2 != nil {
		return fmt.Errorf("error parsing the configuration file: %v", err2)
	}

	return nil
}

/*
verifyServerConfig - This method will verify that the configuration file has
what it needs.
TODO finish fleshing this out
*/
func (c *ServerConfigType) verifyServerConfig() error {
	var problemsFound = 0

	// --------------------------------------------------
	// Global Configuration
	// --------------------------------------------------
	problemsFound += c.verifyGlobalConfig()

	// --------------------------------------------------
	// Global HTML Configuration
	// --------------------------------------------------
	// If HTML output is turned off globally, then there no need to check the
	// configuration and verify everything is present and valid.
	if c.HTML.Enabled.Value == true {
		problemsFound += c.verifyGlobalHTMLConfig()
	} else {
		c.Logger.Infoln("CONFIG: HTML output is disabled in the global configuration")
	}

	// --------------------------------------------------
	// Discovery Server
	// --------------------------------------------------
	// Only verify the Discovery server configuration if it is enabled.
	if c.DiscoveryServer.Enabled == true {
		problemsFound += c.verifyDiscoveryConfig()
	} else {
		c.Logger.Infoln("CONFIG: The Discovery Server is not enabled in the configuration file")
	}

	if c.DiscoveryServer.Enabled == true && c.HTML.Enabled.Value == true {
		problemsFound += c.verifyDiscoveryHTMLConfig()
	} else {
		c.Logger.Infoln("CONFIG: The Discovery Server is enabled in the configuration file but HTML output is not")
	}

	// --------------------------------------------------
	// API Root Server
	// --------------------------------------------------
	// Only verify the API Root server configuration if it is enabled.
	if c.APIRootServer.Enabled == true {
		problemsFound += c.verifyAPIRootConfig()
	} else {
		c.Logger.Infoln("CONFIG: The API Root Server is not enabled in the configuration file")
	}

	if c.APIRootServer.Enabled == true && c.HTML.Enabled.Value == true {
		problemsFound += c.verifyAPIRootHTMLConfig()
	} else {
		c.Logger.Infoln("CONFIG: The API Root Server is enabled in the configuration file but HTML output is not")
	}

	if problemsFound > 0 {
		c.Logger.Println("ERROR: The configuration has", problemsFound, "error(s)")
		return errors.New("ERROR: Configuration errors found")
	}
	return nil
}

/*
exists - This method checks to see if the filename exists on the file system
*/
func (c *ServerConfigType) exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
