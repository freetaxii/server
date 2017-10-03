// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"encoding/json"
	"github.com/freetaxii/libtaxii2/objects/apiRoot"
	"github.com/freetaxii/libtaxii2/objects/collection"
	"github.com/freetaxii/libtaxii2/objects/discovery"
	"github.com/gorilla/mux"
	"log"
	"os"
)

// Log Level 1 = basic system logging information, sent to STDOUT unless Enabled = true then it is logged to a file
// Log Level 2 =
// Log Level 3 = detailed debugging information and code troubleshooting (like key variable changes)
// Log Level 4 =
// Log Level 5 = RAW packet/message decode and output

// ServerConfigType - This type defines the configuration for the entire server.
type ServerConfigType struct {
	Router *mux.Router
	Global struct {
		Protocol        string
		Listen          string
		Prefix          string
		DbConfig        bool
		DbType          string
		DbFile          string
		HTMLTemplateDir string
		TLSKey          string
		TLSCrt          string
	}
	Logging struct {
		Enabled  bool
		LogLevel int
		LogFile  string
	}
	DiscoveryServer struct {
		Enabled          bool
		HTMLEnabled      bool
		HTMLTemplateFile string
		Services         []DiscoveryServiceType
	}
	APIRootServer struct {
		Enabled           bool
		HTMLEnabled       bool
		HTMLTemplateFiles HTMLFilesType
		Services          []APIRootServiceType
	}
	DiscoveryResources  map[string]discovery.DiscoveryType
	APIRootResources    map[string]apiRoot.APIRootType
	CollectionResources map[string]CollectionServiceType
}

// DiscoveryServiceType - This struct represents an instance of a Discovery server.
// If someone tries to set the 'resourcepath' directive in the configuration file it
// will get overwritten in code.
type DiscoveryServiceType struct {
	Enabled          bool
	Name             string
	ResourcePath     string // Set in verifyDiscoveryConfig()
	HTMLEnabled      bool   // Set in verifyDiscoveryHTMLConfig()
	HTMLTemplateFile string
	HTMLTemplatePath string // Set in verifyDiscoveryHTMLConfig() = Prefix + HTMLTemplateDir
	LogLevel         int    // Set in verifyDiscoveryConfig()
	ResourceID       string
}

// APIRootServiceType - This struct represents an instance of an API Root server.
// If someone tries to set the 'path' directive in the configuration file it
// will just get overwritten in code.
type APIRootServiceType struct {
	Enabled           bool
	Name              string
	ResourcePath      string // Set in verifyAPIRootConfig()
	HTMLEnabled       bool   // Set in verifyAPIRootHTMLConfig()
	HTMLTemplateFiles HTMLFilesType
	HTMLTemplatePath  string // Set in verifyAPIRootHTMLConfig() = Prefix + HTMLTemplateDir
	LogLevel          int    // Set in verifyAPIRootConfig()
	ResourceID        string
	Collections       struct {
		Enabled      bool
		ResourcePath string // Set in verifyAPIRootConfig()
		Members      []string
	}
}

type HTMLFilesType struct {
	APIRoot     string
	Collections string
	Collection  string
	Objects     string
}

type CollectionServiceType struct {
	Location     string
	RemoteConfig struct {
		Name string
		URL  string
	}
	Resource collection.CollectionType
}

// --------------------------------------------------
// Load Configuration File and Parse JSON
// --------------------------------------------------

// LoadServerConfig - This methods takes in a string value representing a
// filename of the configuration file and loads the configuration into memory.
func (ezt *ServerConfigType) LoadServerConfig(filename string) {
	// TODO - Need to make make a validation check for the configuration file

	// Open and read configuration file
	sysConfigFileData, err := os.Open(filename)
	defer sysConfigFileData.Close()
	if err != nil {
		log.Fatalf("error opening configuration file: %v", err)
	}

	// --------------------------------------------------
	// Decode JSON configuration file
	// --------------------------------------------------
	// Use decoder instead of unmarshal so we can handle stream data
	decoder := json.NewDecoder(sysConfigFileData)
	err = decoder.Decode(ezt)

	if err != nil {
		log.Fatalf("error parsing configuration file %v", err)
	}

	if ezt.Logging.LogLevel >= 5 {
		log.Println("DEBUG-5: System Configuration Dump")
		log.Printf("%+v\n", ezt)
	}
}
