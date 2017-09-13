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

// -----------------------------------------------------------------------------
// StartObjectsService - This will look to see which objects are part of this collection.
// This method should only be called from the startCollectionService()
// param: apiRootIndex - an integer representing the current location of the API-Root for loop
// param: collection - a string value representing a specific collection in the map
func (this *ServerConfigType) startObjectsService(apiRootIndex int, collection string) {
	collectionid := this.CollectionResources[collection].Resource.Id

	// Make a copy of just the elements that we need to process the request and nothing more.
	// This is done to prevent sending the entire server config in to each handler
	var ts server.ServerHandlerType
	ts.Type = "Objects"
	ts.Path = this.ApiRootService.Services[apiRootIndex].Collections.Path + collectionid + "/" + "objects/"
	ts.Html = this.ApiRootService.Html
	ts.HtmlFile = this.ApiRootService.Services[apiRootIndex].HtmlFiles.Objects
	ts.HtmlPath = this.System.HtmlDir + ts.HtmlFile
	ts.LogLevel = this.Logging.LogLevel
	ts.Location = this.CollectionResources[collection].Location
	ts.RemoteConfig = this.CollectionResources[collection].RemoteConfig

	// --------------------------------------------------
	// Start a Objects handler
	// --------------------------------------------------
	log.Println("Starting TAXII Object service of:", ts.Path)
	if this.CollectionResources[collection].Location == "remote" {
		this.Router.HandleFunc(ts.Path, ts.ObjectsServerRemoteHandler).Methods("GET")
	} else if this.CollectionResources[collection].Location == "local" {
		this.Router.HandleFunc(ts.Path, ts.ObjectsServerHandler).Methods("GET")
	}
}

// -----------------------------------------------------------------------------
// StartObjectByIdService - This will return a specific object that is part of this collection.
// This method should only be called from the startCollectionService()
// param: apiRootIndex - an integer representing the current location of the API-Root for loop
// param: collectionid - a string value representing a specific collection id
func (this *ServerConfigType) startObjectByIdService(apiRootIndex int, collection string) {
	collectionid := this.CollectionResources[collection].Resource.Id

	// Make a copy of just the elements that we need to process the request and nothing more.
	// This is done to prevent sending the entire server config in to each handler
	var ts server.ServerHandlerType
	ts.Type = "ObjectId"
	ts.Path = this.ApiRootService.Services[apiRootIndex].Collections.Path + collectionid + "/" + "objects/" + "{objectid}/"
	ts.Html = this.ApiRootService.Html
	ts.HtmlFile = this.ApiRootService.Services[apiRootIndex].HtmlFiles.Objects
	ts.HtmlPath = this.System.HtmlDir + ts.HtmlFile
	ts.LogLevel = this.Logging.LogLevel
	ts.Location = this.CollectionResources[collection].Location
	ts.RemoteConfig = this.CollectionResources[collection].RemoteConfig

	// --------------------------------------------------
	// Start a Objects handler
	// --------------------------------------------------
	log.Println("Starting TAXII Object By ID service of:", ts.Path)
	if this.CollectionResources[collection].Location == "remote" {
		this.Router.HandleFunc(ts.Path, ts.ObjectsServerRemoteHandler).Methods("GET")
	} else if this.CollectionResources[collection].Location == "local" {
		this.Router.HandleFunc(ts.Path, ts.ObjectsServerHandler).Methods("GET")
	}
}
