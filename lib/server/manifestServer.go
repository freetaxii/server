// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file in the root of the source tree.

package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/freetaxii/freetaxii-server/lib/headers"
	"github.com/freetaxii/libstix2/defs"
	"github.com/freetaxii/libstix2/resources"
	"github.com/gologme/log"
)

/*
ManifestServerHandler - This method will handle all of the requests for STIX
objects from the TAXII server.
*/
func (s *ServerHandlerType) ManifestServerHandler(w http.ResponseWriter, r *http.Request) {
	var mediaType string
	var taxiiHeader headers.HttpHeaderType
	var objectNotFound = false
	var q resources.CollectionQueryType
	var addedFirst, addedLast string

	log.Infoln("INFO: Found Request on the Manifest Server Handler from", r.RemoteAddr, "for collection:", s.CollectionID)

	// If trace is enabled in the logger, than lets decode the HTTP Request and
	// dump it to the logs
	if log.GetLevel("trace") {
		taxiiHeader.DebugHttpRequest(r)
	}

	httpHeaderAccept := r.Header.Get("Accept")
	// httpHeaderRange := r.Header.Get("Range")

	// myregexp := regexp.MustCompile(`^items \d+-\d+$`)
	// if myregexp.MatchString(httpHeaderRange) {
	// 	rangeData := strings.Split(httpHeaderRange, " ")
	// 	if rangeData[0] == "items" {
	// 		values := strings.Split(rangeData[1], "-")
	// 		q.RangeBegin, _ = strconv.Atoi(values[0])
	// 		q.RangeEnd, _ = strconv.Atoi(values[1])

	// 		log.Debugln("DEBUG-3: Client", r.RemoteAddr, "sent the following range parameters:", values[0], values[1])
	// 	}
	// }

	// ----------------------------------------------------------------------
	//
	// Handle URL Parameters
	//
	// ----------------------------------------------------------------------

	urlParameters := r.URL.Query()
	log.Debugln("DEBUG: Client", r.RemoteAddr, "sent the following url parameters:", urlParameters)

	q.CollectionID = s.CollectionID
	errURLParameters := q.ProcessURLParameters(urlParameters)
	if errURLParameters != nil {
		log.Warnln("WARN: invalid URL parameters from client", r.RemoteAddr, "with URL parameters", urlParameters, errURLParameters)
	}

	results, err := s.DS.GetManifestData(q)

	if err != nil {
		taxiiError := resources.NewError()
		title := "ERROR: " + err.Error()
		taxiiError.SetTitle(title)
		desc := "The requested had the following problem: " + err.Error()
		taxiiError.SetDescription(desc)
		taxiiError.SetHTTPStatus("404")
		s.Resource = taxiiError
		objectNotFound = true
		log.Infoln("INFO: Sending error response to", r.RemoteAddr, "due to:", err.Error())

	} else {
		s.Resource = results.ManifestData
		addedFirst = results.DateAddedFirst
		addedLast = results.DateAddedLast
		log.Infoln("INFO: Sending response to", r.RemoteAddr)

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
	w.Header().Add("X-TAXII-Date-Added-First", addedFirst)
	w.Header().Add("X-TAXII-Date-Added-Last", addedLast)
	// contentRangeHeaderValue := "items " + strconv.Itoa(results.RangeBegin) + "-" + strconv.Itoa(results.RangeEnd) + "/" + strconv.Itoa(results.Size)
	// w.Header().Add("Content-Range", contentRangeHeaderValue)

	if strings.Contains(httpHeaderAccept, defs.TAXII_MEDIA_TYPE) {
		mediaType = defs.TAXII_MEDIA_TYPE + "; " + defs.TAXII_VERSION + "; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)

		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusPartialContent)
		}
		j.Encode(s.Resource)

	} else if strings.Contains(httpHeaderAccept, "application/json") {
		mediaType = "application/json; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)

		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusPartialContent)
		}
		j.SetIndent("", "    ")
		j.Encode(s.Resource)

	} else if s.HTMLEnabled == true && strings.Contains(httpHeaderAccept, "text/html") {
		mediaType = "text/html; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
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
			log.Fatal("Unable to create JSON Message")
		}
		s.Resource = string(jsondata)

		// ----------------------------------------------------------------------
		// Setup HTML Template
		// ----------------------------------------------------------------------
		htmlFullPath := s.HTMLTemplatePath + "/" + s.HTMLTemplateFile
		htmlTemplateResource := template.Must(template.ParseFiles(htmlFullPath))
		htmlTemplateResource.ExecuteTemplate(w, s.HTMLTemplateFile, s)

	} else {
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}
}
