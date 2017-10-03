// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"github.com/freetaxii/freetaxii-server/lib/config"
)

// --------------------------------------------------
// Setup Handler Structs
// --------------------------------------------------

type TAXIIServerHandlerType struct {
	// Since we are using a single handler for multiple taxii messages, we need to know the resource type
	Type             string
	ResourcePath     string // This is used in the HTML output
	HTMLEnabled      bool
	HTMLTemplateFile string
	HTMLTemplatePath string
	LogLevel         int
	Resource         interface{}
}

// ServerHandlerType - This struct will handle the discovery, api_root, collections, collection, etc
type ServerHandlerType struct {
	// Since we are using a single handler for multiple taxii messages, we need to know the resource type
	Type string
	// Needed
	HTMLEnabled bool
	// Needed
	HTMLTemplateFile string
	// HTMLPath = Full path + filename
	HTMLTemplatePath string
	// Needed
	LogLevel int
	// Not used
	ResourcePath string
	Resource     interface{}
	Location     string
	RemoteConfig struct {
		Name string
		URL  string
	}
}

// NewDiscoveryHandler - This method will make a copy of the elements found in
// the main configuration for this Discovery Service and copy them here. We do
// this so that we do not send the entire configuration to a handler. Also, this
// enables us to create a generic handler that can fulfill requests for all of
// the TAXII handlers because we can pre-format the data to be in the correct
// from.
func (ezt *TAXIIServerHandlerType) NewDiscoveryHandler(c config.DiscoveryServiceType) {
	ezt.Type = "Discovery"
	ezt.ResourcePath = c.ResourcePath
	ezt.HTMLEnabled = c.HTMLEnabled
	ezt.HTMLTemplateFile = c.HTMLTemplateFile
	ezt.HTMLTemplatePath = c.HTMLTemplatePath
	ezt.LogLevel = c.LogLevel
}

func (ezt *TAXIIServerHandlerType) NewAPIRootHandler(c config.APIRootServiceType) {
	ezt.Type = "API-Root"
	ezt.ResourcePath = c.ResourcePath
	ezt.HTMLEnabled = c.HTMLEnabled
	ezt.HTMLTemplateFile = c.HTMLTemplateFiles.APIRoot
	ezt.HTMLTemplatePath = c.HTMLTemplatePath
	ezt.LogLevel = c.LogLevel
}

func (ezt *TAXIIServerHandlerType) NewCollectionsHandler(c config.APIRootServiceType) {
	ezt.Type = "Collections"
	ezt.ResourcePath = c.Collections.ResourcePath
	ezt.HTMLEnabled = c.HTMLEnabled
	ezt.HTMLTemplateFile = c.HTMLTemplateFiles.Collections
	ezt.HTMLTemplatePath = c.HTMLTemplatePath
	ezt.LogLevel = c.LogLevel
}

func (ezt *TAXIIServerHandlerType) NewCollectionHandler(c config.APIRootServiceType, path string) {
	ezt.Type = "Collection"
	ezt.ResourcePath = path
	ezt.HTMLEnabled = c.HTMLEnabled
	ezt.HTMLTemplateFile = c.HTMLTemplateFiles.Collection
	ezt.HTMLTemplatePath = c.HTMLTemplatePath
	ezt.LogLevel = c.LogLevel
}
