// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"os"
	"strings"

	"github.com/freetaxii/freetaxii-server/lib/config"
	"github.com/freetaxii/libstix2/datastore"
	"github.com/freetaxii/libstix2/resources"
	"github.com/gologme/log"
)

// --------------------------------------------------
// Setup Handler Struct
// --------------------------------------------------

/*
ServerHandlerType - This type will hold the data elements required to process
all TAXII requests.
*/
type ServerHandlerType struct {
	Logger       *log.Logger
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

func New(logger *log.Logger) (ServerHandlerType, error) {
	var s ServerHandlerType

	if logger == nil {
		s.Logger = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		s.Logger = logger
	}

	return s, nil
}

/*
NewDiscoveryHandler - This function will prepare the data for the Discovery handler.
*/
func NewDiscoveryHandler(logger *log.Logger, c config.DiscoveryServiceType, r resources.DiscoveryType) (ServerHandlerType, error) {
	s, _ := New(logger)
	s.URLPath = c.FullPath
	s.HTMLEnabled = c.HTML.Enabled.Value
	s.HTMLTemplate = c.HTML.FullTemplatePath + c.HTML.TemplateFiles.Discovery.Value
	s.Resource = r
	return s, nil
}

/*
NewAPIRootHandler - This function will prepare the data for the API Root handler.
*/
func NewAPIRootHandler(logger *log.Logger, c config.APIRootServiceType, r resources.APIRootType) (ServerHandlerType, error) {
	s, _ := New(logger)
	s.URLPath = c.FullPath
	s.HTMLEnabled = c.HTML.Enabled.Value
	s.HTMLTemplate = c.HTML.FullTemplatePath + c.HTML.TemplateFiles.APIRoot.Value
	s.Resource = r
	return s, nil
}

/*
NewCollectionsHandler - This function will prepare the data for the Collections handler.
*/
func NewCollectionsHandler(logger *log.Logger, c config.APIRootServiceType, r resources.CollectionsType) (ServerHandlerType, error) {
	s, _ := New(logger)
	s.URLPath = c.FullPath + "collections/"
	s.HTMLEnabled = c.HTML.Enabled.Value
	s.HTMLTemplate = c.HTML.FullTemplatePath + c.HTML.TemplateFiles.Collections.Value
	s.Resource = r
	return s, nil
}

/*
NewCollectionHandler - This function will prepare the data for the Collection handler.
*/
func NewCollectionHandler(logger *log.Logger, c config.APIRootServiceType, r resources.CollectionType) (ServerHandlerType, error) {
	s, _ := New(logger)
	s.URLPath = c.FullPath + "collections/" + r.ID + "/"
	s.HTMLEnabled = c.HTML.Enabled.Value
	s.HTMLTemplate = c.HTML.FullTemplatePath + c.HTML.TemplateFiles.Collection.Value
	s.Resource = r
	return s, nil
}

/*
NewObjectsHandler - This function will prepare the data for the Objects handler.
*/
func NewObjectsHandler(logger *log.Logger, c config.APIRootServiceType, collectionID string) (ServerHandlerType, error) {
	s, _ := New(logger)
	s.URLPath = c.FullPath + "collections/" + collectionID + "/objects/"
	s.HTMLEnabled = c.HTML.Enabled.Value
	s.HTMLTemplate = c.HTML.FullTemplatePath + c.HTML.TemplateFiles.Objects.Value
	s.CollectionID = collectionID
	return s, nil
}

/*
NewObjectsByIDHandler - This function will prepare the data for the Objects by ID handler.
*/
func NewObjectsByIDHandler(logger *log.Logger, c config.APIRootServiceType, collectionID string) (ServerHandlerType, error) {
	s, _ := New(logger)
	s.URLPath = c.FullPath + "collections/" + collectionID + "/objects/{objectid}/"
	s.HTMLEnabled = c.HTML.Enabled.Value
	s.HTMLTemplate = c.HTML.FullTemplatePath + c.HTML.TemplateFiles.Objects.Value
	s.CollectionID = collectionID
	return s, nil
}

/*
NewManifestHandler - This function will prepare the data for the Manifest handler.
*/
func NewManifestHandler(logger *log.Logger, c config.APIRootServiceType, collectionID string) (ServerHandlerType, error) {
	s, _ := New(logger)
	s.URLPath = c.FullPath + "collections/" + collectionID + "/manifest/"
	s.HTMLEnabled = c.HTML.Enabled.Value
	s.HTMLTemplate = c.HTML.FullTemplatePath + c.HTML.TemplateFiles.Manifest.Value
	s.CollectionID = collectionID
	return s, nil
}

// ----------------------------------------------------------------------
// Private Methods - ServerHandlerType
// ----------------------------------------------------------------------

/*
processURLParameters - This method will process all of the URL parameters from
an HTTP request.
*/
func (s *ServerHandlerType) processURLParameters(q *resources.CollectionQueryType, values map[string][]string) error {

	if values["id"] != nil {
		q.STIXID = strings.Split(values["id"][0], ",")
	}

	if values["type"] != nil {
		q.STIXType = strings.Split(values["type"][0], ",")
	}

	if values["version"] != nil {
		q.STIXVersion = strings.Split(values["version"][0], ",")
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
