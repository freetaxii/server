// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/freetaxii/freetaxii-server/lib/headers"
	"github.com/freetaxii/libstix2/defs"
	"github.com/gologme/log"
)

/*
TAXIIServerHandler - This method handles all requests for TAXII media type responses
*/
func (s *ServerHandlerType) TAXIIServerHandler(w http.ResponseWriter, r *http.Request) {
	var mediaType string
	var httpHeaderAccept string
	var taxiiHeader headers.HttpHeaderType

	log.Infoln("INFO: Found", s.Type, "request from", r.RemoteAddr, "at", r.RequestURI)

	// If trace is enabled in the logger, than lets decode the HTTP Request and
	// dump it to the logs
	if log.GetLevel("trace") {
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
		j.Encode(s.Resource)

	} else if strings.Contains(httpHeaderAccept, "application/json") {
		mediaType = "application/json; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		w.WriteHeader(http.StatusOK)
		j.SetIndent("", "    ")
		j.Encode(s.Resource)

	} else if s.HTMLEnabled == true && strings.Contains(httpHeaderAccept, "text/html") {
		mediaType = "text/html; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		w.WriteHeader(http.StatusOK)

		// ----------------------------------------------------------------------
		// Setup HTML Template
		// ----------------------------------------------------------------------
		// htmlTemplateResource := template.Must(template.Parse(s.HTMLTemplate))
		// htmlTemplateResource.Execute(w, s)

	} else {
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}

	log.Infoln("INFO: Sending", s.Type, "response to", r.RemoteAddr)
}
