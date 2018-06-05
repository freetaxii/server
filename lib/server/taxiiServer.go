// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/freetaxii/freetaxii-server/lib/headers"
	"github.com/freetaxii/libstix2/defs"
)

/*
DiscoveryHandler - This method will handle all Discovery requests
*/
func (s *ServerHandlerType) DiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("INFO: Found Discovery request from", r.RemoteAddr, "at", r.RequestURI)
	s.baseHandler(w, r)
}

/*
APIRootHandler - This method will handle all API Root requests
*/
func (s *ServerHandlerType) APIRootHandler(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("INFO: Found API Root request from", r.RemoteAddr, "at", r.RequestURI)
	s.baseHandler(w, r)
}

/*
CollectionsHandler - This method will handle all Collections requests
*/
func (s *ServerHandlerType) CollectionsHandler(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("INFO: Found Collections request from", r.RemoteAddr, "at", r.RequestURI)
	s.baseHandler(w, r)
}

/*
CollectionHandler - This method will handle all Collection requests
*/
func (s *ServerHandlerType) CollectionHandler(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("INFO: Found Collection request from", r.RemoteAddr, "at", r.RequestURI)
	s.baseHandler(w, r)
}

/*
baseHandler - This method handles all requests for the following TAXII
media type responses: Discovery, API-Root, Collections, and Collection
*/
func (s *ServerHandlerType) baseHandler(w http.ResponseWriter, r *http.Request) {
	var taxiiHeader headers.HttpHeaderType
	var acceptHeader headers.AcceptHeaderType
	acceptHeader.ParseTAXII(r.Header.Get("Accept"))

	// If trace is enabled in the logger, than decode the HTTP Request to the log
	if s.Logger.GetLevel("trace") {
		taxiiHeader.DebugHttpRequest(r)
	}

	// --------------------------------------------------
	// Encode outgoing response message
	// --------------------------------------------------

	// Setup JSON stream encoder
	j := json.NewEncoder(w)

	// Set header for TLS
	w.Header().Add("Strict-Transport-Security", "max-age=86400; includeSubDomains")

	if acceptHeader.TAXII21 == true {
		w.Header().Set("Content-Type", defs.CONTENT_TYPE_TAXII21)
		w.WriteHeader(http.StatusOK)
		j.Encode(s.Resource)

	} else if acceptHeader.TAXII20 == true {
		w.Header().Set("Content-Type", defs.CONTENT_TYPE_TAXII20)
		w.WriteHeader(http.StatusOK)
		j.Encode(s.Resource)

	} else if acceptHeader.JSON == true {
		w.Header().Set("Content-Type", defs.CONTENT_TYPE_JSON)
		w.WriteHeader(http.StatusOK)
		j.SetIndent("", "    ")
		j.Encode(s.Resource)

	} else if s.HTMLEnabled == true && acceptHeader.HTML == true {
		w.Header().Set("Content-Type", defs.CONTENT_TYPE_HTML)
		w.WriteHeader(http.StatusOK)

		// ----------------------------------------------------------------------
		// Setup HTML Template
		// ----------------------------------------------------------------------
		htmlTemplateResource := template.Must(template.ParseFiles(s.HTMLTemplate))
		htmlTemplateResource.Execute(w, s)

	} else {
		w.Header().Set("Content-Type", defs.CONTENT_TYPE_TAXII21)
		w.WriteHeader(http.StatusNotAcceptable)
	}
}
