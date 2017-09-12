// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"github.com/freetaxii/freetaxii-server/lib/server"
	"log"
)

// StartCollection - This will look to see which collections are defined for this
// Collections group in this API Root. If they are enabled, it will setup handlers for it.
// The HandleFunc passes in copy of the Collection Resource and the extra meta data
// that it needs to process the request.
// This method should only be called from the startCollectionsService()
// param: apiRootIndex - an integer representing the current location of the API-Root for loop
func (this *ServerConfigType) startCollectionService(apiRootIndex int) {

	// We need to loop through each collection in this API Root and setup a handler for it.
	for i, value := range this.ApiRootService.Services[apiRootIndex].Collections.Resource.Collections {

		// Make a copy of just the elements that we need to process the request and nothing more.
		// This is done to prevent sending the entire server config in to each handler
		var ts server.ServerHandlerType
		ts.Type = "Collection"
		ts.Path = this.ApiRootService.Services[apiRootIndex].Collections.Path + value.Id + "/"
		ts.HtmlFile = this.ApiRootService.Services[apiRootIndex].Collection.HtmlFile
		ts.HtmlPath = this.System.HtmlDir + ts.HtmlFile
		ts.LogLevel = this.Logging.LogLevel
		ts.Resource = this.ApiRootService.Services[apiRootIndex].Collections.Resource.Collections[i]

		// --------------------------------------------------
		// Start a Collection handler
		// --------------------------------------------------
		log.Println("Starting TAXII Collection service of:", ts.Path)

		// We do not need to check to see if the collection is enabled and readable/writeable because that was already done
		// TODO add support for post if the colleciton is writeable
		this.Router.HandleFunc(ts.Path, ts.TaxiiServerHandler).Methods("GET")

		// --------------------------------------------------
		// Start a Objects handler
		// --------------------------------------------------
		// TODO you need to know which collection this object is in, so you only show the right objects.
		var taxiiObjects server.ServerHandlerType
		taxiiObjects.Type = "Objects"
		taxiiObjects.Path = ts.Path + "objects/"
		taxiiObjects.HtmlFile = "objectsResource.html"
		taxiiObjects.HtmlPath = this.System.HtmlDir + taxiiObjects.HtmlFile
		taxiiObjects.LogLevel = this.Logging.LogLevel

		log.Println("Starting TAXII Object service of:", taxiiObjects.Path)
		this.Router.HandleFunc(taxiiObjects.Path, taxiiObjects.ObjectsServerHandler).Methods("GET")

		// --------------------------------------------------
		// Start a Object by ID handler
		// --------------------------------------------------
		var taxiiObjectId server.ServerHandlerType
		taxiiObjectId.Type = "Object-ID"
		taxiiObjectId.Path = taxiiObjects.Path + "{objectid}/"
		taxiiObjectId.HtmlFile = "objectsResource.html"
		taxiiObjectId.HtmlPath = this.System.HtmlDir + taxiiObjectId.HtmlFile
		taxiiObjects.LogLevel = this.Logging.LogLevel

		log.Println("Starting TAXII Object by ID service of:", taxiiObjectId.Path)
		this.Router.HandleFunc(taxiiObjectId.Path, taxiiObjectId.ObjectsServerHandler).Methods("GET")
	}
}
