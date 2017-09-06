// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"encoding/json"
	"github.com/freetaxii/libtaxii2/objects/api_root"
	"github.com/freetaxii/libtaxii2/objects/common"
	"github.com/freetaxii/libtaxii2/objects/discovery"
	"log"
	"os"
)

// Log Level 1 = basic system logging information, sent to STDOUT unless Enabled = true then it is logged to a file
// Log Level 2 =
// Log Level 3 = detailed debugging information and code troubleshooting (like key variable changes)
// Log Level 4 =
// Log Level 5 = RAW packet/message decode and output

type ServerType struct {
	System struct {
		Listen          string
		Prefix          string
		DbFile          string
		DbFileFullPath  string
		HtmlDir         string
		HtmlDirFullPath string
	}
	Logging struct {
		Enabled         bool
		LogLevel        int
		LogFile         string
		LogFileFullPath string
	}
	DiscoveryService struct {
		Enabled   bool
		Resources []DiscoveryResourceType
	}
	ApiRootService struct {
		Enabled   bool
		Resources []ApiRootResourceType
	}
}

type DiscoveryResourceType struct {
	common.DiscoveryMetadataProperties
	Resource discovery.DiscoveryType
}

type ApiRootResourceType struct {
	common.ApiRootMetadataProperties
	Resource api_root.ApiRootType
}

// --------------------------------------------------
// Load Configuration File and Parse JSON
// --------------------------------------------------

// LoadServerConfig - This methods takes in one parameter
// param: s - a string value representing a filename of the configuration file
func (this *ServerType) LoadServerConfig(filename string) {
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

	// Lets assign the full paths to a few variables so we can use them later
	this.System.DbFileFullPath = this.System.Prefix + "/" + this.System.DbFile
	this.System.HtmlDirFullPath = this.System.Prefix + "/" + this.System.HtmlDir
	this.Logging.LogFileFullPath = this.System.Prefix + "/" + this.Logging.LogFile

	if this.Logging.LogLevel >= 5 {
		log.Printf("DEBUG-5: System Configuration Dump %+v\n", this)
	}
}
