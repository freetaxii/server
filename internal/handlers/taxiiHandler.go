// Copyright 2015-2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/freetaxii/libstix2/defs"
	"github.com/freetaxii/server/internal/headers"
)

/*
DiscoveryHandler - This method will handle all Discovery requests
*/
func (s *ServerHandler) DiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("INFO: Found Discovery request from", r.RemoteAddr, "at", r.RequestURI)
	s.baseHandler(w, r)
}

/*
APIRootHandler - This method will handle all API Root requests
*/
func (s *ServerHandler) APIRootHandler(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("INFO: Found API Root request from", r.RemoteAddr, "at", r.RequestURI)
	s.baseHandler(w, r)
}

/*
CollectionsHandler - This method will handle all Collections requests
*/
func (s *ServerHandler) CollectionsHandler(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("INFO: Found Collections request from", r.RemoteAddr, "at", r.RequestURI)
	s.baseHandler(w, r)
}

/*
CollectionHandler - This method will handle all Collection requests
*/
func (s *ServerHandler) CollectionHandler(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("INFO: Found Collection request from", r.RemoteAddr, "at", r.RequestURI)
	s.baseHandler(w, r)
}

/*
baseHandler - This method handles all requests for the following TAXII
media type responses: Discovery, API-Root, Collections, and Collection
*/
func (s *ServerHandler) baseHandler(w http.ResponseWriter, r *http.Request) {

	// If trace is enabled in the logger, than decode the HTTP Request to the log
	if s.Logger.GetLevel("trace") {
		headers.DebugHttpRequest(r)
	}

	// --------------------------------------------------
	// 1st Check Authentication
	// --------------------------------------------------

	// If authentication is required and the client does not provide credentials
	// or their credentials do not match, then send an error message.
	// We need to return right here as to prevent further processing.
	if s.Authenticated == true {
		s.Logger.Debugln("DEBUG: Authentication Enabled")
		if s.BasicAuth == true {
			s.Logger.Debugln("DEBUG: Basic Authentication Enabled")
			w.Header().Set("WWW-Authenticate", `Basic realm="Authentication Required"`)
			if success := s.authenticate(r.BasicAuth()); success != true {
				s.Logger.Debugln("DEBUG: Authentication failed for", r.RemoteAddr, "at", r.RequestURI)
				s.sendUnauthenticatedError(w)
				return
			}
		} else {
			// If authentication is enabled, but basic is not, then fail since
			// no other authentication is currently enabled.
			s.Logger.Debugln("DEBUG: Authentication method from", r.RemoteAddr, "at", r.RequestURI, "not supported")
			s.sendUnauthenticatedError(w)
			return
		}
	} // End Authentication Check

	// --------------------------------------------------
	// Check Accept Header Media Type
	// --------------------------------------------------
	var acceptHeader headers.MediaType
	acceptHeader.ParseTAXII(r.Header.Get("Accept"))

	// --------------------------------------------------
	// Encode outgoing response message
	// --------------------------------------------------

	// Set header for TLS
	w.Header().Add("Strict-Transport-Security", "max-age=86400; includeSubDomains")

	if acceptHeader.TAXII21 == true {
		// Setup JSON stream encoder
		j := json.NewEncoder(w)
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_TAXII21)
		w.WriteHeader(http.StatusOK)
		j.Encode(s.Resource)

	} else if acceptHeader.JSON == true {
		// Setup JSON stream encoder
		j := json.NewEncoder(w)
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_JSON)
		w.WriteHeader(http.StatusOK)
		j.SetIndent("", "    ")
		j.Encode(s.Resource)

	} else if s.HTMLEnabled == true && acceptHeader.HTML == true {
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_HTML)
		w.WriteHeader(http.StatusOK)

		// ----------------------------------------------------------------------
		// Setup HTML Template
		// ----------------------------------------------------------------------
		htmlTemplateResource := template.Must(template.ParseFiles(s.HTMLTemplate))
		htmlTemplateResource.Execute(w, s)

	} else {
		s.sendNotAcceptableError(w)
		return
	}
}
