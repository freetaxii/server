// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"github.com/freetaxii/freetaxii-server/lib/config"
	"github.com/freetaxii/libstix2/datastore"
)

// --------------------------------------------------
// Setup Handler Structs
// --------------------------------------------------

// TAXIIServerHandlerType - This type will hold the data elements required to
// process all TAXII media type requests. Since we are using a single handler
// for multiple taxii messages, we need to know the resource type.
type TAXIIServerHandlerType struct {
	Type             string // Used in log messages
	ResourcePath     string // This is used in the HTML output
	HTMLEnabled      bool
	HTMLTemplateFile string
	HTMLTemplatePath string // Prefix + HTMLTemplateDir
	LogLevel         int
	Resource         interface{}
}

// STIXServerHandlerType - This type will hold the data elements required to
// process all STIX media type requests. Since we are using a single handler
// for multiple stix messages, we need to know the resource type.
type STIXServerHandlerType struct {
	ResourcePath     string // This is used in the HTML output
	HTMLEnabled      bool
	HTMLTemplateFile string
	HTMLTemplatePath string // Prefix + HTMLTemplateDir
	LogLevel         int
	CollectionID     string
	RangeMax         int
	DS               datastore.Datastorer
	Resource         interface{}
}

// These methods will copy the elements found in the main configuration file.
// We do this so that we do not send the entire configuration to a handler.
// Also, this enables us to create a generic handler that can fulfill requests
// for all of the TAXII and STIX handlers because we can pre-format the data to
// be in a consistent and correct from.

// NewDiscoveryHandler - This method will make a copy of the elements found in
// the main configuration for this Discovery Service and copy them here.
func (ezt *TAXIIServerHandlerType) NewDiscoveryHandler(c config.DiscoveryServiceType) {
	ezt.Type = "Discovery"
	ezt.ResourcePath = c.ResourcePath
	ezt.HTMLEnabled = c.HTMLEnabled
	ezt.HTMLTemplateFile = c.HTMLBranding.Discovery
	ezt.HTMLTemplatePath = c.HTMLTemplatePath
	ezt.LogLevel = c.LogLevel
}

// NewAPIRootHandler - This method will make a copy of the elements found in
// the main configuration for this API Root Service and copy them here.
func (ezt *TAXIIServerHandlerType) NewAPIRootHandler(c config.APIRootServiceType) {
	ezt.Type = "API-Root"
	ezt.ResourcePath = c.ResourcePath
	ezt.HTMLEnabled = c.HTMLEnabled
	ezt.HTMLTemplateFile = c.HTMLBranding.APIRoot
	ezt.HTMLTemplatePath = c.HTMLTemplatePath
	ezt.LogLevel = c.LogLevel
}

// NewCollectionsHandler - This method will make a copy of the elements found in
// the main configuration for this Collections Service and copy them here.
func (ezt *TAXIIServerHandlerType) NewCollectionsHandler(c config.APIRootServiceType) {
	ezt.Type = "Collections"
	ezt.ResourcePath = c.Collections.ResourcePath
	ezt.HTMLEnabled = c.HTMLEnabled
	ezt.HTMLTemplateFile = c.HTMLBranding.Collections
	ezt.HTMLTemplatePath = c.HTMLTemplatePath
	ezt.LogLevel = c.LogLevel
}

// NewCollectionHandler - This method will make a copy of the elements found in
// the main configuration for this Collection Service and copy them here.
func (ezt *TAXIIServerHandlerType) NewCollectionHandler(c config.APIRootServiceType, path string) {
	ezt.Type = "Collection"
	ezt.ResourcePath = path
	ezt.HTMLEnabled = c.HTMLEnabled
	ezt.HTMLTemplateFile = c.HTMLBranding.Collection
	ezt.HTMLTemplatePath = c.HTMLTemplatePath
	ezt.LogLevel = c.LogLevel
}
