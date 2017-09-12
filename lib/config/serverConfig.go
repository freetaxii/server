// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"encoding/json"
	"github.com/freetaxii/libtaxii2/objects/api_root"
	"github.com/freetaxii/libtaxii2/objects/collection"
	"github.com/freetaxii/libtaxii2/objects/collections"
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

type ServerConfigType struct {
	Router *mux.Router
	System struct {
		Protocol string
		Listen   string
		Prefix   string
		DbConfig bool
		DbFile   string
		HtmlDir  string
		TlsKey   string
		TlsCrt   string
	}
	Logging struct {
		Enabled  bool
		LogLevel int
		LogFile  string
	}
	DiscoveryService struct {
		Enabled  bool
		HtmlFile string
		Services []DiscoveryServiceType
	}
	ApiRootService struct {
		Enabled  bool
		HtmlFile string
		Services []ApiRootServiceType
	}
	AllCollections map[string]CollectionServiceType
}

// If someone tries to set the 'path' directive in the configuration file it will just get overwritten in code.
type DiscoveryServiceType struct {
	Enabled  bool
	Name     string
	Path     string
	HtmlFile string
	Resource discovery.DiscoveryType
}

// If someone tries to set the 'path' directive in the configuration file it will just get overwritten in code.
type ApiRootServiceType struct {
	Enabled     bool
	Name        string
	Path        string
	HtmlFile    string
	Collections struct {
		Enabled          bool
		Path             string
		HtmlFile         string
		ValidCollections []string
		Resource         collections.CollectionsType
	}
	Collection struct {
		HtmlFile string
	}
	Resource api_root.ApiRootType
}

// If someone tries to set the 'path' directive in the configuration file it will just get overwritten in code.
type CollectionServiceType struct {
	Enabled  bool
	Name     string
	Path     string
	Resource collection.CollectionType
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
	defer sysConfigFileData.Close()
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

// VerifyServerConfig - This method will verify that the configuration file has what it needs
// TODO finish fleshing this out
func (this *ServerConfigType) VerifyServerConfig() error {
	var err error
	err = this.verifyConfigDirectives()
	if err != nil {
		return err
	}

	err = this.verifyHtmlTemplateFiles()
	if err != nil {
		return err
	}
	return nil
}
