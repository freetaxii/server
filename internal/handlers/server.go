// Copyright 2015-2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package handlers

import (
	"os"
	"strings"

	"github.com/freetaxii/libstix2/datastore"
	"github.com/freetaxii/libstix2/resources/apiroot"
	"github.com/freetaxii/libstix2/resources/collections"
	"github.com/freetaxii/libstix2/resources/discovery"
	"github.com/freetaxii/server/internal/config"
	"github.com/gologme/log"
)

// --------------------------------------------------
// Setup Handler Struct
// --------------------------------------------------

/*
ServerHandler - This type will hold the data elements required to process
all TAXII requests.
*/
type ServerHandler struct {
	Logger            *log.Logger
	URLPath           string // Used in HTML output and to build the URL for the next resource.
	HTMLEnabled       bool   // Is HTML output enabled for this service
	HTMLTemplate      string // The full file path (prefix + HTML template directory + template filename)
	CollectionID      string // The collection ID that is being used
	ServerRecordLimit int    // The maximum number of records that the server will respond with.
	Authenticated     bool   // Is this handler to be authenticated
	BasicAuth         bool   // Is Basic Auth used
	DS                datastore.Datastorer
	Resource          interface{} // This holds the actual resource and is populated in the main freetaxii.go
}

// ----------------------------------------------------------------------
// These methods will copy the elements found in the main configuration file.
// We do this so that we do not send the entire configuration to a handler.
// Also, this enables us to create a generic handler that can fulfill requests
// for all of the TAXII and STIX handlers because we can pre-format the data to
// be in a consistent and correct from.
// ----------------------------------------------------------------------

func New(logger *log.Logger) (ServerHandler, error) {
	var s ServerHandler

	if logger == nil {
		s.Logger = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		s.Logger = logger
	}

	// TODO for right now lets just force this until we can plumb this in through the configuration file
	s.Authenticated = false
	s.BasicAuth = false

	return s, nil
}

/*
NewDiscoveryHandler - This function will prepare the data for the Discovery handler.
*/
func NewDiscoveryHandler(logger *log.Logger, c config.DiscoveryService, r discovery.Discovery) (ServerHandler, error) {
	s, _ := New(logger)
	s.URLPath = c.Path
	s.HTMLEnabled = c.HTML.Enabled.Value
	s.HTMLTemplate = c.HTML.FullTemplatePath + c.HTML.TemplateFiles.Discovery.Value
	s.Resource = r
	return s, nil
}

/*
NewAPIRootHandler - This function will prepare the data for the API Root handler.
*/
func NewAPIRootHandler(logger *log.Logger, api config.APIRootService, r apiroot.APIRoot) (ServerHandler, error) {
	s, _ := New(logger)
	s.URLPath = api.Path
	s.HTMLEnabled = api.HTML.Enabled.Value
	s.HTMLTemplate = api.HTML.FullTemplatePath + api.HTML.TemplateFiles.APIRoot.Value
	s.Resource = r
	return s, nil
}

/*
NewCollectionsHandler - This function will prepare the data for the Collections handler.
*/
func NewCollectionsHandler(logger *log.Logger, api config.APIRootService, r collections.Collections, limit int) (ServerHandler, error) {
	s, _ := New(logger)
	s.URLPath = api.Path + "collections/"
	s.HTMLEnabled = api.HTML.Enabled.Value
	s.HTMLTemplate = api.HTML.FullTemplatePath + api.HTML.TemplateFiles.Collections.Value
	s.Resource = r
	s.ServerRecordLimit = limit
	return s, nil
}

/*
NewCollectionHandler - This function will prepare the data for the Collection handler.
*/
func NewCollectionHandler(logger *log.Logger, api config.APIRootService, r collections.Collection, limit int) (ServerHandler, error) {
	s, _ := New(logger)
	s.URLPath = api.Path + "collections/" + r.ID + "/"
	s.HTMLEnabled = api.HTML.Enabled.Value
	s.HTMLTemplate = api.HTML.FullTemplatePath + api.HTML.TemplateFiles.Collection.Value
	s.Resource = r
	s.ServerRecordLimit = limit
	return s, nil
}

/*
NewObjectsHandler - This function will prepare the data for the Objects handler.
*/
func NewObjectsHandler(logger *log.Logger, api config.APIRootService, collectionID string, limit int) (ServerHandler, error) {
	s, _ := New(logger)
	s.URLPath = api.Path + "collections/" + collectionID + "/objects/"
	s.HTMLEnabled = api.HTML.Enabled.Value
	s.HTMLTemplate = api.HTML.FullTemplatePath + api.HTML.TemplateFiles.Objects.Value
	s.CollectionID = collectionID
	s.ServerRecordLimit = limit
	return s, nil
}

/*
NewObjectsByIDHandler - This function will prepare the data for the Objects by ID handler.
*/
func NewObjectsByIDHandler(logger *log.Logger, api config.APIRootService, collectionID string, limit int) (ServerHandler, error) {
	s, _ := New(logger)
	s.URLPath = api.Path + "collections/" + collectionID + "/objects/{objectid}/"
	s.HTMLEnabled = api.HTML.Enabled.Value
	s.HTMLTemplate = api.HTML.FullTemplatePath + api.HTML.TemplateFiles.Objects.Value
	s.CollectionID = collectionID
	s.ServerRecordLimit = limit
	return s, nil
}

/*
NewManifestHandler - This function will prepare the data for the Manifest handler.
*/
func NewManifestHandler(logger *log.Logger, api config.APIRootService, collectionID string, limit int) (ServerHandler, error) {
	s, _ := New(logger)
	s.URLPath = api.Path + "collections/" + collectionID + "/manifest/"
	s.HTMLEnabled = api.HTML.Enabled.Value
	s.HTMLTemplate = api.HTML.FullTemplatePath + api.HTML.TemplateFiles.Manifest.Value
	s.CollectionID = collectionID
	s.ServerRecordLimit = limit
	return s, nil
}

// ----------------------------------------------------------------------
// Private Methods - ServerHandler
// ----------------------------------------------------------------------

/*
processURLParameters - This method will process all of the URL parameters from
an HTTP request.
*/
func (s *ServerHandler) processURLParameters(q *collections.CollectionQuery, values map[string][]string) error {

	if values["match[id]"] != nil {
		q.STIXID = strings.Split(values["match[id]"][0], ",")
	}

	if values["match[type]"] != nil {
		q.STIXType = strings.Split(values["match[type]"][0], ",")
	}

	if values["match[version]"] != nil {
		q.STIXVersion = strings.Split(values["match[version]"][0], ",")
	}

	if values["added_after"] != nil {
		q.AddedAfter = strings.Split(values["added_after"][0], ",")
	}

	if values["added_before"] != nil {
		q.AddedBefore = strings.Split(values["added_before"][0], ",")
	}

	if values["limit"] != nil {
		q.Limit = strings.Split(values["limit"][0], ",")
	}

	s.Logger.Debugln("DEBUG: URL Parameter ID", q.STIXID)
	s.Logger.Debugln("DEBUG: URL Parameter Type", q.STIXType)
	s.Logger.Debugln("DEBUG: URL Parameter Version", q.STIXVersion)
	s.Logger.Debugln("DEBUG: URL Parameter Added After", q.AddedAfter)
	s.Logger.Debugln("DEBUG: URL Parameter Added Before", q.AddedBefore)
	s.Logger.Debugln("DEBUG: URL Parameter Limit", q.Limit)
	return nil
}
