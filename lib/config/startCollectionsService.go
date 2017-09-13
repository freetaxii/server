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

// --------------------------------------------------------------------------------
// StartCollections - This will look to see if the Collections service is enabled
// in the config file for a given API Root. If it is, it will setup handlers for it.
// The HandleFunc passes in copy of the Collections Resource and the extra meta data
// that it needs to process the request.
// This method should only be called from the StartApiRootService()
// param: apiRootIndex - an integer representing the current location of the API-Root for loop
func (this *ServerConfigType) startCollectionsService(apiRootIndex int) {

	this.ApiRootService.Services[apiRootIndex].Collections.Path = this.ApiRootService.Services[apiRootIndex].Path + "collections/"

	// Make a copy of just the elements that we need to process the request and nothing more.
	// This is done to prevent sending the entire server config in to each handler
	var ts server.ServerHandlerType
	ts.Type = "Collections"
	ts.Html = this.ApiRootService.Html
	ts.Path = this.ApiRootService.Services[apiRootIndex].Collections.Path
	ts.HtmlFile = this.ApiRootService.Services[apiRootIndex].HtmlFiles.Collections
	ts.HtmlPath = this.System.HtmlDir + ts.HtmlFile
	ts.LogLevel = this.Logging.LogLevel
	// This next value will be set down below
	//ts.Resource

	// We need to look in to this instance of the API Root and find out which collections are tied to it
	// Then we can use that ID to pull from the collections list and add them to this list of valid collections

	collections := objects.NewCollections()
	for _, value := range this.ApiRootService.Services[apiRootIndex].Collections.Members {

		// If enabled, only add the collection to the list if the collection can either be read or written to
		if this.CollectionResources[value].Resource.Can_read == true || this.CollectionResources[value].Resource.Can_write == true {
			collections.AddCollection(this.CollectionResources[value].Resource)
		}

	}
	ts.Resource = collections

	log.Println("Starting TAXII Collections service of:", ts.Path)
	this.Router.HandleFunc(ts.Path, ts.TaxiiServerHandler).Methods("GET")

	// --------------------------------------------------
	// Start a Collection handler
	// --------------------------------------------------
	this.startCollectionService(apiRootIndex)

}

// --------------------------------------------------------------------------------
// StartCollection - This will look to see which collections are defined for this
// Collections group in this API Root. If they are enabled, it will setup handlers for it.
// The HandleFunc passes in copy of the Collection Resource and the extra meta data
// that it needs to process the request.
// This method should only be called from the startCollectionsService()
// param: apiRootIndex - an integer representing the current location of the API-Root for loop
func (this *ServerConfigType) startCollectionService(apiRootIndex int) {

	// We need to loop through each collection in this API Root and setup a handler for it.
	for _, value := range this.ApiRootService.Services[apiRootIndex].Collections.Members {

		// Make a copy of just the elements that we need to process the request and nothing more.
		// This is done to prevent sending the entire server config in to each handler
		var ts server.ServerHandlerType
		ts.Type = "Collection"
		ts.Path = this.ApiRootService.Services[apiRootIndex].Collections.Path + this.CollectionResources[value].Resource.Id + "/"
		ts.Html = this.ApiRootService.Html
		ts.HtmlFile = this.ApiRootService.Services[apiRootIndex].HtmlFiles.Collection
		ts.HtmlPath = this.System.HtmlDir + ts.HtmlFile
		ts.LogLevel = this.Logging.LogLevel
		ts.Resource = this.CollectionResources[value].Resource

		// --------------------------------------------------
		// Start a Collection handler
		// --------------------------------------------------
		log.Println("Starting TAXII Collection service of:", ts.Path)

		// We do not need to check to see if the collection is enabled and readable/writeable because that was already done
		// TODO add support for post if the colleciton is writeable
		this.Router.HandleFunc(ts.Path, ts.TaxiiServerHandler).Methods("GET")

		// --------------------------------------------------
		// Start an Objects handler
		// --------------------------------------------------
		// This will pass in the map name not the UUIDv4 collection ID
		this.startObjectsService(apiRootIndex, value)
		this.startObjectByIdService(apiRootIndex, value)
	}
}
