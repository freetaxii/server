// Copyright 2015-2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/freetaxii/libstix2/resources/apiroot"
	"github.com/freetaxii/libstix2/resources/collections"
	"github.com/freetaxii/libstix2/resources/discovery"
	"github.com/gologme/log"
	"github.com/gorilla/mux"
)

/*
ServerConfig - This type defines the configuration for the entire server.
*/
type ServerConfig struct {
	Router *mux.Router
	Logger *log.Logger
	Global struct {
		Prefix            string
		Listen            string
		Protocol          string
		TLSDir            string
		TLSKey            string
		TLSCrt            string
		DbConfig          bool
		DbType            string
		DbFile            string
		ServerRecordLimit int
	}
	HTML struct {
		HTMLConfig
	}
	Logging struct {
		Enabled bool
		Level   int
		LogFile string
	}
	DiscoveryServer struct {
		Enabled  bool
		Services []DiscoveryService
	} `json:"discovery_server,omitempty"`
	APIRootServer struct {
		Enabled  bool
		Services []APIRootService
	} `json:"apiroot_server,omitempty"`
	DiscoveryResources  map[string]discovery.Discovery    `json:"discovery_resources,omitempty"`  // The key in the map is the ResourceID
	APIRootResources    map[string]apiroot.APIRoot        `json:"apiroot_resources,omitempty"`    // The key in the map is the ResourceID
	CollectionResources map[string]collections.Collection `json:"collection_resources,omitempty"` // The key in the map is the ResourceID
}

/*
BaseService - This struct represents the common properties between the
Discovery and API-Root services.

Path          - The URL path for this service
Enabled       - Is this service enabled
ResourceID    - A unique ID for the resource that this service is using
ResourcePath  - The actual full URL path for the resource, used for the handler to know where to listen.
HTML          - The configuration for generating HTML output
*/
type BaseService struct {
	Enabled    bool       // User defined in configuration file
	Path       string     // User defined in configuration file
	ResourceID string     // User defined in configuration file
	HTML       HTMLConfig // User defined in configuration file or set in the verify scripts.
}

/*
DiscoveryService - This struct represents an instance of a Discovery server.
If someone tries to set the 'resourcepath' directive in the configuration file it
will get overwritten in code.
*/
type DiscoveryService struct {
	BaseService
}

/*
APIRootService - This struct represents an instance of an API Root server.
If someone tries to set the 'path' directive in the configuration file it
will just get overwritten in code.
ReadAccess - This is a list of collection resource IDs that may have GET access
at the API Root level
WriteAccess - This is a list of collection resource IDs that may have POST access
at the API Root level
*/
type APIRootService struct {
	BaseService
	Collections struct {
		Enabled     bool     // User defined in configuration file
		ReadAccess  []string // User defined in configuration file.
		WriteAccess []string // User defined in configuration file.
	}
}

/*
HTMLConfig - This struct holds the configuration elements for generating HTML
output. This is used at the top level of the configuration file as well as in
each individual service. This means individual services can have a different
HTML configuration. I needed to setup my own types for JSON boolean and strings
since leaving it blank at a child level, would have equaled "false" or "". This
would have been equivalent to turning it off, which is not what is wanted. Leaving
it blank would mean to inherit from the parent. But since Go is a strictly typed
language, you need to create a type that can handle that case.

Enabled           - Is HTML enabled for this service
TemplateDir       - The location of the template files relative to the base of the application (prefix)
FullTemplatePath  - The full path of the template directory (prefix + TemplateDir)
TemplateFiles     - The HTML template filenames in the template directory for the following services
*/
type HTMLConfig struct {
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
New - This function will load the current configuration from a file, verify that
the configuration is correct, and then return a ServerConfig type.
*/
func New(logger *log.Logger, filename string) (ServerConfig, error) {
	var c ServerConfig
	var err error

	if logger == nil {
		c.Logger = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		c.Logger = logger
	}

	err = c.loadServerConfig(filename)
	if err != nil {
		return c, err
	}

	// In addition to checking the configuration for completeness the verify
	// process will also populate some of the helper values.
	err = c.Verify()
	if err != nil {
		return c, err
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
func (c *ServerConfig) loadServerConfig(filename string) error {
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
exists - This method checks to see if the filename exists on the file system.
This is used by several of the configuration directive checks, basically anytime
there is a filename defined in the configuration file, this is called to check
to see if that file actually exists on the file system.
*/
func (c *ServerConfig) exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
