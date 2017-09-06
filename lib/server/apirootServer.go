// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"github.com/freetaxii/freetaxii-server/defs"
	"github.com/freetaxii/freetaxii-server/lib/headers"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// ApiRootServerHandler - This method takes in three parameters. The last parameter
// the index is so that this handler will know which directory service is being called
// in case there is more than one.
// param: w - http.ResponseWriter
// param: r - *http.Request
// param: index - An integer that lets this method know which discovery service is being handled by this instance
func (this *ServerType) ApiRootServerHandler(w http.ResponseWriter, r *http.Request, index int) {
	var mediaType string
	var httpHeaderAccept string
	var jsondata []byte
	var formatpretty bool = false
	var taxiiHeader headers.HttpHeaderType

	// Setup HTML template
	var htmlResourceFile string = "apirootResource.html"
	var htmlResource string = this.System.HtmlDir + "/" + htmlResourceFile
	var htmlTemplateResource = template.Must(template.ParseFiles(htmlResource))

	if this.Logging.LogLevel >= 3 {
		log.Printf("DEBUG-3: Found Request on API Root Server Handler from %s", r.RemoteAddr)
	}

	// We need to put this first so that during debugging we can see problems
	// that will generate errors below.
	if this.Logging.LogLevel >= 5 {
		taxiiHeader.DebugHttpRequest(r)
	}

	// --------------------------------------------------
	// Decode incoming request message
	// --------------------------------------------------
	httpHeaderAccept = r.Header.Get("Accept")

	if strings.Contains(httpHeaderAccept, defs.TAXII_MEDIA_TYPE) {
		mediaType = defs.TAXII_MEDIA_TYPE + "; " + defs.TAXII_VERSION + "; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		formatpretty = false
		jsondata = this.createApiRootResponse(formatpretty, index)
		w.Write(jsondata)
	} else if strings.Contains(httpHeaderAccept, "text/html") {
		mediaType = "text/html; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		htmlTemplateResource.ExecuteTemplate(w, htmlResourceFile, this.ApiRootService.Resources[index])
	} else {
		mediaType = "application/json; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		formatpretty = true
		jsondata = this.createApiRootResponse(formatpretty, index)
		w.Write(jsondata)
	}

	if this.Logging.LogLevel >= 1 {
		log.Println("DEBUG-1: Sending API Root Response to", r.RemoteAddr)
	}
}
