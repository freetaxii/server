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

// StartDiscoveryService - This will look to see if there are any Discovery services
// defined in the config file. If there are, it will loop through the list and setup
// handlers for each one of them. The HandleFunc passes in copy of the Discovery Resource
// and the extra meta data that it needs to process the request.
// retval: serviceCounter - an integer that keeps track of how many services were started
func (this *ServerConfigType) StartDiscoveryService() int {
	var serviceCounter int = 0
	for index := range this.DiscoveryService.Services {

		// Check to see if this entry is actually enabled
		if this.DiscoveryService.Services[index].Enabled == true {

			// Make a copy of just the elements that we need to process the request and nothing more.
			// This is done to prevent sending the entire server config in to each handler
			var ts server.ServerHandlerType
			ts.Type = "Discovery"
			ts.Path = this.DiscoveryService.Services[index].Path
			ts.Html = this.DiscoveryService.Html
			ts.HtmlFile = this.DiscoveryService.Services[index].HtmlFile
			ts.HtmlPath = this.System.HtmlDir + ts.HtmlFile
			ts.LogLevel = this.Logging.LogLevel
			ts.Resource = this.DiscoveryResources[this.DiscoveryService.Services[index].Resource]

			log.Println("Starting TAXII Discovery service at:", ts.Path)
			this.Router.HandleFunc(ts.Path, ts.TaxiiServerHandler).Methods("GET")
			serviceCounter++
		}
	}
	return serviceCounter
}
