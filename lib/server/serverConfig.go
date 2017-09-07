// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"encoding/json"
	"github.com/freetaxii/libtaxii2/objects/api_root"
	"github.com/freetaxii/libtaxii2/objects/collection"
	"github.com/freetaxii/libtaxii2/objects/discovery"
	"log"
	"os"
)

// Log Level 1 = basic system logging information, sent to STDOUT unless Enabled = true then it is logged to a file
// Log Level 2 =
// Log Level 3 = detailed debugging information and code troubleshooting (like key variable changes)
// Log Level 4 =
// Log Level 5 = RAW packet/message decode and output

type ServerConfigType struct {
	System struct {
		Protocol string
		Listen   string
		Prefix   string
		DbConfig bool
		DbFile   string
		HtmlDir  string
	}
	Logging struct {
		Enabled  bool
		LogLevel int
		LogFile  string
	}
	DiscoveryService struct {
		Enabled  bool
		Services []DiscoveryServiceType
	}
	ApiRootService struct {
		Enabled  bool
		Services []ApiRootServiceType
	}
	Collections map[string]CollectionServiceType
}

type DiscoveryServiceType struct {
	Enabled  bool
	Path     string
	Resource discovery.DiscoveryType
}

type ApiRootServiceType struct {
	Enabled     bool
	Path        string
	Collections []string
	Resource    api_root.ApiRootType
}

type CollectionServiceType struct {
	Enabled  bool
	Resource collection.CollectionType
}

// --------------------------------------------------
// Setup Handler Structs
// --------------------------------------------------
// This struct will handle discovery, api_root, collections, collection, etc
type ServerHandlerType struct {
	HtmlDir  string
	LogLevel int
	Path     string
	Resource interface{}
}

// --------------------------------------------------
// Load Configuration File and Parse JSON
// --------------------------------------------------

// LoadServerConfig - This methods takes in one parameter
// param: s - a string value representing a filename of the configuration file
func (this *ServerConfigType) LoadServerConfig(filename string) {
	// TODO - Need to make make a validation check for the configuration file

	// Open and read configuration file
	sysConfigFileData, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error opening configuration file: %v", err)
	}

	// --------------------------------------------------
	// Decode JSON configuration file
	// --------------------------------------------------
	// Use decoder instead of unmarshal so we can handle stream data
	decoder := json.NewDecoder(sysConfigFileData)
	err = decoder.Decode(this)

	if err != nil {
		log.Fatalf("error parsing configuration file %v", err)
	}

	if this.Logging.LogLevel >= 5 {
		log.Printf("DEBUG-5: System Configuration Dump %+v\n", this)
	}
}
