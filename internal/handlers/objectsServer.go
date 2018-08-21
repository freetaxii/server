// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/freetaxii/libstix2/defs"
	"github.com/freetaxii/libstix2/objects"
	"github.com/freetaxii/libstix2/resources"
	"github.com/freetaxii/libstix2/stixid"
	"github.com/freetaxii/server/internal/headers"
	"github.com/gorilla/mux"
)

/*
ObjectsServerHandler - This method will handle all of the requests for STIX
objects from the TAXII server.
*/
func (s *ServerHandler) ObjectsServerHandler(w http.ResponseWriter, r *http.Request) {
	var taxiiHeader headers.HttpHeader
	var acceptHeader headers.MediaType
	acceptHeader.ParseSTIX(r.Header.Get("Accept"))

	var objectNotFound = false
	var addedFirst, addedLast string
	q := resources.NewCollectionQuery(s.CollectionID, s.ServerRecordLimit)

	s.Logger.Infoln("INFO: Found GET Request on the Objects Server Handler from", r.RemoteAddr, "for collection:", s.CollectionID)

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

	// 		s.Logger.Debugln("DEBUG: Client", r.RemoteAddr, "sent the following range parameters:", values[0], values[1])
	// 	}
	// }

	// ----------------------------------------------------------------------
	// Handle URL Parameters
	// ----------------------------------------------------------------------

	urlvars := mux.Vars(r)
	if urlvars["objectid"] != "" {
		urlObjectID := urlvars["objectid"]
		s.Logger.Debugln("DEBUG: Client", r.RemoteAddr, "sent URL path value:", urlObjectID)

		// TODO check to see if objectid is valid first, change to make work with custom objects
		if stixid.ValidSTIXID(urlObjectID) {
			q.STIXID = append(q.STIXID, urlObjectID)
		} else if stixid.ValidSTIXObjectType(urlObjectID) {
			q.STIXType = append(q.STIXType, urlObjectID)
		}
	}

	urlParameters := r.URL.Query()
	s.Logger.Debugln("DEBUG: Client", r.RemoteAddr, "sent URL parameters:", urlParameters)

	errURLParameters := s.processURLParameters(q, urlParameters)
	if errURLParameters != nil {
		s.Logger.Warnln("WARN: invalid URL parameters from client", r.RemoteAddr, "with URL parameters", urlParameters, errURLParameters)
	}

	results, err := s.DS.GetBundle(*q)

	if err != nil {
		taxiiError := resources.NewError()
		title := "ERROR: " + err.Error()
		taxiiError.SetTitle(title)
		desc := "The requested had the following problem: " + err.Error()
		taxiiError.SetDescription(desc)
		taxiiError.SetHTTPStatus("404")
		s.Resource = taxiiError
		objectNotFound = true
		s.Logger.Infoln("INFO: Sending error response to", r.RemoteAddr, "due to:", err.Error())

	} else {
		s.Resource = results.BundleData
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

	if acceptHeader.STIX21 == true {
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_STIX21)

		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusPartialContent)
		}
		j.Encode(s.Resource)

	} else if acceptHeader.STIX20 == true {
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_STIX20)

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

/*
ObjectsServerWriteHandler - This method will handle all POST requests of STIX
objects from the TAXII server.
*/
func (s *ServerHandler) ObjectsServerWriteHandler(w http.ResponseWriter, r *http.Request) {
	var taxiiHeader headers.HttpHeader
	var acceptHeader headers.MediaType
	var contentHeader headers.MediaType
	acceptHeader.ParseSTIX(r.Header.Get("Accept"))
	contentHeader.ParseSTIX(r.Header.Get("Content-type"))

	s.Logger.Infoln("INFO: Found POST Request on the Objects Server from", r.RemoteAddr, "for collection:", s.CollectionID)

	// If trace is enabled in the logger, than decode the HTTP Request to the log
	if s.Logger.GetLevel("trace") {
		taxiiHeader.DebugHttpRequest(r)
	}

	// ----------------------------------------------------------------------
	// Decode the bundle object itself, but leave the objects array as an
	// array of raw JSON object objects, we will decode each one later.
	// ----------------------------------------------------------------------
	b, err := objects.DecodeBundle(r.Body)
	if err != nil {
		s.Logger.Warnln("WARN: Could not decode provided bundle")

		// TODO if this is an error we need to eject right here and sent error message back to client.
	}

	// TODO first check content-type header

	// ----------------------------------------------------------------------
	// Decode each object in the bundle one at a time. If the object is valid
	// write it off to the datastore.
	// Lets keep a count of the number of objects that are successful and the
	// number that are not successful in addition to a total count
	// ----------------------------------------------------------------------
	totalCount := 0
	successCount := 0
	failureCount := 0
	for _, v := range b.Objects {
		totalCount++
		s.Logger.Debugln("DEBUG: Processing bundle object number", totalCount)

		// First, decode the first object from the bundle if it succeeds try to
		// add it to the datastore
		o, id, err := objects.DecodeObject(v)
		if err != nil {
			// TODO Track something to send error back to client in status resource
			s.Logger.Warnln("WARN: Error decoding object in bundle", err)
			failureCount++
			// If there is an error, lets just skip and move on to the next object
			continue
		}

		// Add the object to the datastore, if the decode was successful
		s.Logger.Debugln("DEBUG: Adding object", id, "to the datastore")
		err = s.DS.AddObject(o)
		if err != nil {
			s.Logger.Warnln("WARN: Error adding object", id, "to datastore", err)
			failureCount++
			// If there was an error, lets just skip and move on to the next object
			continue
		}
		successCount++

		// If the add was successful then lets add an entry in to the collection
		// record table.
		entry := resources.CreateCollectionRecord(s.CollectionID, id)
		s.Logger.Debugln("DEBUG: Adding Collection Entry of", s.CollectionID, id)
		err = s.DS.AddTAXIIObject(entry)
		if err != nil {
			s.Logger.Debugln(err)
		}
	}
	s.Logger.Debugln("DEBUG: Total number of objects in Bundle", totalCount)
	s.Logger.Debugln("DEBUG: Total objects successfully added to datastore", successCount)
	s.Logger.Debugln("DEBUG: Total objects that failed to be added to the datastore", failureCount)
	// unmarshal content and write data

	//results, err := s.DS.GetBundle(*q)

	// if err != nil {
	// 	taxiiError := resources.NewError()
	// 	title := "ERROR: " + err.Error()
	// 	taxiiError.SetTitle(title)
	// 	desc := "The requested had the following problem: " + err.Error()
	// 	taxiiError.SetDescription(desc)
	// 	taxiiError.SetHTTPStatus("404")
	// 	s.Resource = taxiiError
	// 	objectNotFound = true
	// 	s.Logger.Infoln("INFO: Sending error response to", r.RemoteAddr, "due to:", err.Error())

	// } else {

	s.Logger.Infoln("INFO: Sending response to", r.RemoteAddr)
	// }

	// --------------------------------------------------
	// Encode outgoing response message
	// --------------------------------------------------

	// Setup JSON stream encoder for response
	j := json.NewEncoder(w)

	// Set header for TLS
	w.Header().Add("Strict-Transport-Security", "max-age=86400; includeSubDomains")

	if acceptHeader.TAXII21 == true {
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_TAXII21)
		w.WriteHeader(http.StatusAccepted)
		j.Encode(s.Resource)

	} else if acceptHeader.TAXII20 == true {
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_TAXII20)
		w.WriteHeader(http.StatusAccepted)
		j.Encode(s.Resource)

	} else if acceptHeader.JSON == true {
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_JSON)
		w.WriteHeader(http.StatusAccepted)
		j.SetIndent("", "    ")
		j.Encode(s.Resource)

	} else if s.HTMLEnabled == true && acceptHeader.HTML == true {
		w.Header().Set("Content-Type", defs.MEDIA_TYPE_HTML)
		w.WriteHeader(http.StatusAccepted)

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
