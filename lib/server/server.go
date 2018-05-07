// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"github.com/freetaxii/freetaxii-server/lib/config"
	"github.com/freetaxii/libstix2/datastore"
	"github.com/freetaxii/libstix2/resources"
)

// --------------------------------------------------
// Setup Handler Struct
// --------------------------------------------------

/*
ServerHandlerType - This type will hold the data elements required to process
all TAXII requests.
*/
type ServerHandlerType struct {
	URLPath      string // Used in HTML output and to build the URL for the next resource.
	HTMLEnabled  bool   // Is HTML output enabled for this service
	HTMLTemplate string // The full file path (prefix + HTML template directory + template filename)
	CollectionID string
	DS           datastore.Datastorer
	Resource     interface{} // This holds the actual resource and is populated in the main freetaxii.go
}

// ----------------------------------------------------------------------
// These methods will copy the elements found in the main configuration file.
// We do this so that we do not send the entire configuration to a handler.
// Also, this enables us to create a generic handler that can fulfill requests
// for all of the TAXII and STIX handlers because we can pre-format the data to
// be in a consistent and correct from.
// ----------------------------------------------------------------------

/*
NewDiscoveryHandler - This function will prepare the data for the Discovery handler.
*/
func NewDiscoveryHandler(c config.DiscoveryServiceType, r resources.DiscoveryType) (ServerHandlerType, error) {
	var s ServerHandlerType
	s.URLPath = c.FullPath
	s.HTMLEnabled = c.HTML.Enabled.Value
	s.HTMLTemplate = c.HTML.FullTemplatePath + c.HTML.TemplateFiles.Discovery.Value
	s.Resource = r
	return s, nil
}

/*
NewAPIRootHandler - This function will prepare the data for the API Root handler.
*/
func NewAPIRootHandler(c config.APIRootServiceType) (ServerHandlerType, error) {
	var s ServerHandlerType
	s.URLPath = c.FullPath
	s.HTMLEnabled = c.HTML.Enabled.Value
	s.HTMLTemplate = c.HTML.FullTemplatePath + c.HTML.TemplateFiles.APIRoot.Value
	return s, nil
}

/*
NewCollectionsHandler - This function will prepare the data for the Collections handler.
*/
func NewCollectionsHandler(c config.APIRootServiceType) (ServerHandlerType, error) {
	var s ServerHandlerType
	s.URLPath = c.Collections.FullPath
	s.HTMLEnabled = c.HTML.Enabled.Value
	s.HTMLTemplate = c.HTML.FullTemplatePath + c.HTML.TemplateFiles.Collections.Value
	return s, nil
}

/*
NewCollectionHandler - This function will prepare the data for the Collection handler.
*/
func NewCollectionHandler(c config.APIRootServiceType, path string) (ServerHandlerType, error) {
	var s ServerHandlerType
	s.URLPath = path
	s.HTMLEnabled = c.HTML.Enabled.Value
	s.HTMLTemplate = c.HTML.FullTemplatePath + c.HTML.TemplateFiles.Collection.Value
	return s, nil
}
