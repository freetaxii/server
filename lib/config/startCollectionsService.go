// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"github.com/freetaxii/freetaxii-server/lib/server"
	"github.com/freetaxii/libtaxii2/objects"
	"log"
)

// StartCollections - This will look to see if the Collections service is enabled
// in the config file for a given API Root. If it is, it will setup handlers for it.
// The HandleFunc passes in copy of the Collections Resource and the extra meta data
// that it needs to process the request.
// This method should only be called from the StartApiRootService()
// param: apiRootIndex - an integer representing the current location of the API-Root for loop
func (this *ServerConfigType) startCollectionsService(apiRootIndex int) {

	// Make a copy of just the elements that we need to process the request and nothing more.
	// This is done to prevent sending the entire server config in to each handler
	var ts server.ServerHandlerType
	ts.Type = "Collections"
	ts.Path = this.ApiRootService.Services[apiRootIndex].Collections.Path
	ts.HtmlFile = this.ApiRootService.Services[apiRootIndex].Collections.HtmlFile
	ts.HtmlPath = this.System.HtmlDir + ts.HtmlFile
	ts.LogLevel = this.Logging.LogLevel
	// This next value will be set down below
	//ts.Resource

	// We need to look in to this instance of the API Root and find out which collections are tied to it
	// Then we can use that ID to pull from the collections list and add them to this list of valid collections

	collections := objects.NewCollections()
	for _, value := range this.ApiRootService.Services[apiRootIndex].Collections.ValidCollections {

		// Only add the collection if it is enabled
		if this.AllCollections[value].Enabled == true {

			// If enabled, only add the collection to the list if the collection can either be read or written to
			if this.AllCollections[value].Resource.Can_read == true || this.AllCollections[value].Resource.Can_write == true {
				collections.AddCollection(this.AllCollections[value].Resource)
			}
		}

	}
	ts.Resource = collections
	this.ApiRootService.Services[apiRootIndex].Collections.Resource = collections

	log.Println("Starting TAXII Collections service of:", ts.Path)
	this.Router.HandleFunc(ts.Path, ts.TaxiiServerHandler).Methods("GET")

	// --------------------------------------------------
	// Start a Collection handler
	// --------------------------------------------------
	this.startCollectionService(apiRootIndex)

}
