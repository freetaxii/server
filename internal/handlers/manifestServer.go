// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file in the root of the source tree.

package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/freetaxii/libstix2/defs"
	"github.com/freetaxii/libstix2/resources/collections"
	"github.com/freetaxii/libstix2/resources/taxiierror"
	"github.com/freetaxii/server/internal/headers"
)

/*
ManifestServerHandler - This method will handle all of the requests for STIX
objects from the TAXII server.
*/
func (s *ServerHandler) ManifestServerHandler(w http.ResponseWriter, r *http.Request) {
	var taxiiHeader headers.HttpHeader
	var acceptHeader headers.MediaType

	// --------------------------------------------------
	// 1st check authentication
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
	}

	acceptHeader.ParseTAXII(r.Header.Get("Accept"))

	var objectNotFound = false
	var addedFirst, addedLast string
	q := collections.NewCollectionQuery(s.CollectionID, s.ServerRecordLimit)

	s.Logger.Infoln("INFO: Found Request on the Manifest Server Handler from", r.RemoteAddr, "for collection:", s.CollectionID)

	// If trace is enabled in the logger, than decode the HTTP Request to the log
	if s.Logger.GetLevel("trace") {
		taxiiHeader.DebugHttpRequest(r)
	}

	// httpHeaderRange := r.Header.Get("Range")

	// myregexp := regexp.MustCompile(`^items \d+-\d+$`)
	// if myregexp.MatchString(httpHeaderRange) {
	// 	rangeData := strings.Split(httpHeaderRange, " ")
	// 	if rangeData[0] == "items" {
	// 		values := strings.Split(rangeData[1], "-")
	// 		q.RangeBegin, _ = strconv.Atoi(values[0])
	// 		q.RangeEnd, _ = strconv.Atoi(values[1])

	// 		s.Logger.Debugln("DEBUG-3: Client", r.RemoteAddr, "sent the following range parameters:", values[0], values[1])
	// 	}
	// }

	// ----------------------------------------------------------------------
	// Handle URL Parameters
	// ----------------------------------------------------------------------
	urlParameters := r.URL.Query()

	if len(urlParameters) > 0 {
		s.Logger.Debugln("DEBUG: Client", r.RemoteAddr, "sent the following (", len(urlParameters), ") url parameters:", urlParameters)

		errURLParameters := s.processURLParameters(q, urlParameters)
		if errURLParameters != nil {
			s.Logger.Warnln("WARN: invalid URL parameters from client", r.RemoteAddr, "with URL parameters", urlParameters, errURLParameters)
		}
	}

	results, err := s.DS.GetManifestData(*q)

	if err != nil {
		taxiiError := taxiierror.New()
		title := "ERROR: " + err.Error()
		taxiiError.SetTitle(title)
		desc := "The requested had the following problem: " + err.Error()
		taxiiError.SetDescription(desc)
		taxiiError.SetHTTPStatus("404")
		s.Resource = taxiiError
		objectNotFound = true
		s.Logger.Infoln("INFO: Sending error response to", r.RemoteAddr, "due to:", err.Error())

	} else {
		s.Resource = results.ManifestData
		addedFirst = results.DateAddedFirst
		addedLast = results.DateAddedLast
		s.Logger.Infoln("INFO: Sending response to", r.RemoteAddr)

	}

	// --------------------------------------------------
	// Encode outgoing response message
	// --------------------------------------------------

	// Setup JSON stream encoder
	j := json.NewEncoder(w)

	// Set header for TLS
	w.Header().Add("Strict-Transport-Security", "max-age=86400; includeSubDomains")
	w.Header().Add("X-TAXII-Date-Added-First", addedFirst)
	w.Header().Add("X-TAXII-Date-Added-Last", addedLast)
	// contentRangeHeaderValue := "items " + strconv.Itoa(results.RangeBegin) + "-" + strconv.Itoa(results.RangeEnd) + "/" + strconv.Itoa(results.Size)
	// w.Header().Add("Content-Range", contentRangeHeaderValue)

	if acceptHeader.TAXII21 == true {
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_TAXII21)

		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusPartialContent)
		}
		j.Encode(s.Resource)

	} else if acceptHeader.TAXII20 == true {
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_TAXII20)

		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusPartialContent)
		}
		j.Encode(s.Resource)

	} else if acceptHeader.JSON == true {
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_JSON)

		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusPartialContent)
		}
		j.SetIndent("", "    ")
		j.Encode(s.Resource)

	} else if s.HTMLEnabled == true && acceptHeader.HTML == true {
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_HTML)
		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusPartialContent)
		}

		// I needed to convert this to actual JSON since if I just used
		// s.Resource like in other handlers I would get the string output of
		// a Golang struct which is not the same. The reason it works else where
		// is I am not printing the whole object, but rather, referencing the
		// parts as I need them.
		jsondata, err := json.MarshalIndent(s.Resource, "", "    ")
		if err != nil {
			s.Logger.Fatal("Unable to create JSON Message")
		}
		s.Resource = string(jsondata)

		// ----------------------------------------------------------------------
		// Setup HTML Template
		// ----------------------------------------------------------------------
		htmlTemplateResource := template.Must(template.ParseFiles(s.HTMLTemplate))
		htmlTemplateResource.Execute(w, s)

	} else {
		w.WriteHeader(http.StatusNotAcceptable)
	}
}
