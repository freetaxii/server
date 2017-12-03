// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"encoding/json"
	"github.com/freetaxii/freetaxii-server/lib/headers"
	"github.com/freetaxii/libstix2/defs"
	"html/template"
	"log"
	"net/http"
	"strings"
)

/*
TAXIIServerHandler - This method takes in two parameters and handles all
requests for TAXII media type responses
*/
func (ezt *TAXIIServerHandlerType) TAXIIServerHandler(w http.ResponseWriter, r *http.Request) {
	var mediaType string
	var httpHeaderAccept string
	var taxiiHeader headers.HttpHeaderType

	if ezt.LogLevel >= 3 {
		log.Println("DEBUG-3: Found Request on the", ezt.Type, "Server Handler from", r.RemoteAddr)
	}

	// We need to put this first so that during debugging we can see problems
	// that will generate errors below.
	if ezt.LogLevel >= 5 {
		taxiiHeader.DebugHttpRequest(r)
	}

	// --------------------------------------------------
	//
	// Encode outgoing response message
	//
	// --------------------------------------------------

	// Setup JSON stream encoder
	j := json.NewEncoder(w)

	// Set header for TLS
	w.Header().Add("Strict-Transport-Security", "max-age=86400; includeSubDomains")

	httpHeaderAccept = r.Header.Get("Accept")

	if strings.Contains(httpHeaderAccept, defs.TAXII_MEDIA_TYPE) {
		mediaType = defs.TAXII_MEDIA_TYPE + "; " + defs.TAXII_VERSION + "; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		w.WriteHeader(http.StatusOK)
		j.Encode(ezt.Resource)

	} else if strings.Contains(httpHeaderAccept, "application/json") {
		mediaType = "application/json; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		w.WriteHeader(http.StatusOK)
		j.SetIndent("", "    ")
		j.Encode(ezt.Resource)

	} else if ezt.HTMLEnabled == true && strings.Contains(httpHeaderAccept, "text/html") {
		mediaType = "text/html; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		w.WriteHeader(http.StatusOK)

		// ----------------------------------------------------------------------
		// Setup HTML Template
		// ----------------------------------------------------------------------
		htmlFullPath := ezt.HTMLTemplatePath + "/" + ezt.HTMLTemplateFile
		htmlTemplateResource := template.Must(template.ParseFiles(htmlFullPath))
		htmlTemplateResource.ExecuteTemplate(w, ezt.HTMLTemplateFile, ezt)

	} else {
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}

	if ezt.LogLevel >= 3 {
		log.Println("DEBUG-3: Sending", ezt.Type, "Response to", r.RemoteAddr)
	}
}
