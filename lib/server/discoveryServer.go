// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"encoding/json"
	"github.com/freetaxii/freetaxii-server/defs"
	"github.com/freetaxii/freetaxii-server/lib/headers"
	"github.com/freetaxii/libtaxii2/objects"
	"log"
	"net/http"
	"strings"
)

func (this *ServerType) DiscoveryServerHandler(w http.ResponseWriter, r *http.Request, index int) {
	var mediaType string
	var httpHeaderAccept string
	var jsondata []byte
	var formatpretty bool = false
	var taxiiHeader headers.HttpHeaderType

	if this.SysConfig.Logging.LogLevel >= 3 {
		log.Printf("DEBUG-3: Found Request on Discovery Server Handler from %s", r.RemoteAddr)
	}

	// We need to put this first so that during debugging we can see problems
	// that will generate errors below.
	if this.SysConfig.Logging.LogLevel >= 5 {
		taxiiHeader.DebugHttpRequest(r)
	}

	// --------------------------------------------------
	// Decode incoming request message
	// --------------------------------------------------
	httpHeaderAccept = r.Header.Get("Accept")

	if strings.Contains(httpHeaderAccept, defs.TAXII_MEDIA_TYPE) {
		mediaType = defs.TAXII_MEDIA_TYPE + "; charset=utf-8"
		formatpretty = false
	} else if strings.Contains(httpHeaderAccept, "text/html") {
		mediaType = "text/html; charset=utf-8"
		formatpretty = true
	} else {
		mediaType = "application/json; charset=utf-8"
		formatpretty = true
	}
	jsondata = this.createDiscoveryResponse(formatpretty, index)
	w.Header().Set("Content-Type", mediaType)
	w.Write(jsondata)
	if this.SysConfig.Logging.LogLevel >= 1 {
		log.Println("DEBUG-1: Sending Discovery Response to", r.RemoteAddr)
	}
}

// --------------------------------------------------
// Create a TAXII Discovery Response Message
// --------------------------------------------------

func (this *ServerType) createDiscoveryResponse(formatPretty bool, index int) []byte {
	var jsondata []byte
	var err error
	tm := objects.NewDiscovery()

	// TODO pull these from the database or the ServerType object if it has been loaded from database
	tm.SetTitle(this.DiscoveryService.Resources[index].Title)
	tm.SetDescription(this.DiscoveryService.Resources[index].Description)
	tm.SetContact(this.DiscoveryService.Resources[index].Contact)
	tm.SetDefault(this.DiscoveryService.Resources[index].Default)
	for _, apiroot := range this.DiscoveryService.Resources[index].Api_roots {
		tm.AddApiRoot(apiroot)
	}

	if formatPretty == true {
		jsondata, err = json.MarshalIndent(tm, "", "    ")
	} else {
		jsondata, err = json.Marshal(tm)
	}

	if err != nil {
		// If we can not create a status message then there is something
		// wrong with the APIs and nothing is going to work.
		log.Fatal("Unable to create Discovery Response Message")
	}
	return jsondata
}
