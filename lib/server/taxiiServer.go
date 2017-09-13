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

// -----------------------------------------------------------------------------
// TaxiiServerHandler - This method takes in two parameters and handles all requests for TAXII media type responses
// param: w - http.ResponseWriter
// param: r - *http.Request
func (this *ServerHandlerType) TaxiiServerHandler(w http.ResponseWriter, r *http.Request) {
	var mediaType string
	var httpHeaderAccept string
	var jsondata []byte
	var formatpretty bool = false
	var taxiiHeader headers.HttpHeaderType

	// Setup HTML template
	var htmlTemplateResource = template.Must(template.ParseFiles(this.HtmlPath))

	if this.LogLevel >= 3 {
		log.Println("DEBUG-3: Found Request on the", this.Type, "Server Handler from", r.RemoteAddr)
	}

	// We need to put this first so that during debugging we can see problems
	// that will generate errors below.
	if this.LogLevel >= 5 {
		taxiiHeader.DebugHttpRequest(r)
	}

	w.Header().Add("Strict-Transport-Security", "max-age=86400; includeSubDomains")

	// --------------------------------------------------
	// Decode incoming request message
	// --------------------------------------------------
	httpHeaderAccept = r.Header.Get("Accept")

	if strings.Contains(httpHeaderAccept, defs.TAXII_MEDIA_TYPE) {
		mediaType = defs.TAXII_MEDIA_TYPE + "; " + defs.TAXII_VERSION + "; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		formatpretty = false
		jsondata = this.createTAXIIResponse(formatpretty)
		w.WriteHeader(http.StatusOK)
		w.Write(jsondata)
	} else if strings.Contains(httpHeaderAccept, "application/json") {
		mediaType = "application/json; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		formatpretty = true
		jsondata = this.createTAXIIResponse(formatpretty)
		w.WriteHeader(http.StatusOK)
		w.Write(jsondata)
	} else if this.Html == true && strings.Contains(httpHeaderAccept, "text/html") {
		mediaType = "text/html; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		w.WriteHeader(http.StatusOK)
		htmlTemplateResource.ExecuteTemplate(w, this.HtmlFile, this)
	} else {
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}

	if this.LogLevel >= 3 {
		log.Println("DEBUG-3: Sending", this.Type, "Response to", r.RemoteAddr)
	}
}
