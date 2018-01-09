// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"encoding/json"
	"github.com/freetaxii/freetaxii-server/lib/headers"
	"github.com/freetaxii/libstix2/datastore"
	"github.com/freetaxii/libstix2/defs"
	"github.com/freetaxii/libstix2/resources"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

/*
ManifestServerHandler - This method will handle all of the requests for STIX
objects from the TAXII server.
*/
func (ezt *STIXServerHandlerType) ManifestServerHandler(w http.ResponseWriter, r *http.Request) {
	var mediaType string
	var taxiiHeader headers.HttpHeaderType
	var objectNotFound = false
	var q datastore.QueryType
	var addedFirst, addedLast string

	if ezt.LogLevel >= 3 {
		log.Println("DEBUG-3: Found Request on the Manifest Server Handler from", r.RemoteAddr, "for collection:", ezt.CollectionID)
	}

	if ezt.LogLevel >= 5 {
		taxiiHeader.DebugHttpRequest(r)
	}

	httpHeaderAccept := r.Header.Get("Accept")
	httpHeaderRange := r.Header.Get("Range")

	myregexp := regexp.MustCompile(`^items \d+-\d+$`)
	if myregexp.MatchString(httpHeaderRange) {
		rangeData := strings.Split(httpHeaderRange, " ")
		if rangeData[0] == "items" {
			values := strings.Split(rangeData[1], "-")
			q.RangeBegin, _ = strconv.Atoi(values[0])
			q.RangeEnd, _ = strconv.Atoi(values[1])

			if ezt.LogLevel >= 3 {
				log.Println("DEBUG-3: Client", r.RemoteAddr, "sent the following range parameters:", values[0], values[1])
			}
		}
	}

	urlParameters := r.URL.Query()
	if ezt.LogLevel >= 3 {
		log.Println("DEBUG-3: Client", r.RemoteAddr, "sent the following url parameters:", urlParameters)
	}

	q.CollectionID = ezt.CollectionID

	if urlParameters["match[id]"] != nil {
		q.STIXID = urlParameters["match[id]"]
	}

	if urlParameters["match[type]"] != nil {
		q.STIXType = urlParameters["match[type]"]
	}

	if urlParameters["match[version]"] != nil {
		q.STIXVersion = urlParameters["match[version]"]
	}

	if urlParameters["added_after"] != nil {
		q.AddedAfter = urlParameters["added_after"]
	}

	q.RangeMax = ezt.RangeMax

	objectsInCollection, metaData, err := ezt.DS.GetManifestFromCollection(q)

	if err != nil {
		taxiiError := resources.NewError()
		title := "ERROR: " + err.Error()
		taxiiError.SetTitle(title)
		desc := "The requested had the following problem: " + err.Error()
		taxiiError.SetDescription(desc)
		taxiiError.SetHTTPStatus("404")
		ezt.Resource = taxiiError
		objectNotFound = true
		if ezt.LogLevel >= 3 {
			log.Println("DEBUG-3: Sending error response to", r.RemoteAddr, "due to:", err.Error())
		}
	} else {
		ezt.Resource = *objectsInCollection
		addedFirst = metaData.DateAddedFirst
		addedLast = metaData.DateAddedLast
		if ezt.LogLevel >= 3 {
			log.Println("DEBUG-3: Sending response to", r.RemoteAddr)
		}
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
	contentRangeHeaderValue := "items " + strconv.Itoa(metaData.RangeBegin) + "-" + strconv.Itoa(metaData.RangeEnd) + "/" + strconv.Itoa(metaData.Size)
	w.Header().Add("Content-Range", contentRangeHeaderValue)

	if strings.Contains(httpHeaderAccept, defs.TAXII_MEDIA_TYPE) {
		mediaType = defs.TAXII_MEDIA_TYPE + "; " + defs.TAXII_VERSION + "; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)

		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusPartialContent)
		}
		j.Encode(ezt.Resource)

	} else if strings.Contains(httpHeaderAccept, "application/json") {
		mediaType = "application/json; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)

		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusPartialContent)
		}
		j.SetIndent("", "    ")
		j.Encode(ezt.Resource)

	} else if ezt.HTMLEnabled == true && strings.Contains(httpHeaderAccept, "text/html") {
		mediaType = "text/html; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusPartialContent)
		}

		// I needed to convert this to actual JSON since if I just used
		// ezt.Resource like in other handlers I would get the string output of
		// a Golang struct which is not the same. The reason it works else where
		// is I am not printing the whole object, but rather, referencing the
		// parts as I need them.
		jsondata, err := json.MarshalIndent(ezt.Resource, "", "    ")
		if err != nil {
			log.Fatal("Unable to create JSON Message")
		}
		ezt.Resource = string(jsondata)

		// ----------------------------------------------------------------------
		// Setup HTML Template
		// ----------------------------------------------------------------------
		htmlFullPath := ezt.HTMLTemplatePath + "/" + ezt.HTMLTemplateFile
		htmlTemplateResource := template.Must(template.ParseFiles(htmlFullPath))
		htmlTemplateResource.ExecuteTemplate(w, ezt.HTMLTemplateFile, ezt)

	} else {
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}
}
