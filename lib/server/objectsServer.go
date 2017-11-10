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
	"github.com/freetaxii/libstix2/objects"
	"github.com/freetaxii/libstix2/resources"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// ObjectsServerHandler - This method will handle all of the requests for STIX
// objects from the TAXII server.
func (ezt *STIXServerHandlerType) ObjectsServerHandler(w http.ResponseWriter, r *http.Request) {
	var mediaType string
	var httpHeaderAccept string
	var taxiiHeader headers.HttpHeaderType
	var objectNotFound = false

	urlvars := mux.Vars(r)
	urlObjectID := urlvars["objectid"]

	// Setup a STIX Bundle to be used for response
	stixBundle := objects.NewBundle()

	// ----------------------------------------------------------------------
	// Setup HTML Template - only if HTMLEnabled is true
	// ----------------------------------------------------------------------
	var htmlTemplateResource *template.Template
	if ezt.HTMLEnabled == true {
		var htmlFullPath = ezt.HTMLTemplatePath + "/" + ezt.HTMLTemplateFile
		htmlTemplateResource = template.Must(template.ParseFiles(htmlFullPath))
	}

	if ezt.LogLevel >= 3 {
		log.Println("DEBUG-3: Found Request on the", ezt.Type, "Server Handler from", r.RemoteAddr)
	}

	// We need to put this first so that during debugging we can see problems
	// that will generate errors below.
	if ezt.LogLevel >= 5 {
		taxiiHeader.DebugHttpRequest(r)
	}

	// Is this a request for a specific object ID /objects/{objectid}?
	if urlObjectID == "" {
		// Get a list of objects that are in the collection
		allObjects := ezt.DS.GetObjectsInCollection(ezt.CollectionID)
		for _, stixid := range allObjects {
			i, _ := ezt.DS.GetObject(stixid)
			stixBundle.AddObject(i)
		}
		// Add resource to object so we can pass it in to the JSON processor
		ezt.Resource = stixBundle
	} else {
		// If we are looking for just a single object do this part of the if statement
		// TODO make sure this object is in the collection first.
		i, err := ezt.DS.GetObject(urlObjectID)
		if err != nil {
			taxiiError := resources.NewError()
			title := "ERROR: " + err.Error()
			taxiiError.SetTitle(title)
			desc := "The following requested object resource does not exist: " + urlObjectID
			taxiiError.SetDescription(desc)
			taxiiError.SetHTTPStatus("404")
			ezt.Resource = taxiiError
			objectNotFound = true
		} else {
			stixBundle.AddObject(i)
			// Add resource to object so we can pass it in to the JSON processor
			ezt.Resource = stixBundle
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

	httpHeaderAccept = r.Header.Get("Accept")

	if strings.Contains(httpHeaderAccept, defs.STIX_MEDIA_TYPE) {
		mediaType = defs.STIX_MEDIA_TYPE + "; " + defs.STIX_VERSION + "; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)

		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		j.Encode(ezt.Resource)

	} else if strings.Contains(httpHeaderAccept, "application/json") {
		mediaType = "application/json; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)

		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		j.SetIndent("", "    ")
		j.Encode(ezt.Resource)

	} else if ezt.HTMLEnabled == true && strings.Contains(httpHeaderAccept, "text/html") {
		mediaType = "text/html; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		if objectNotFound == true {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
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
		htmlTemplateResource.ExecuteTemplate(w, ezt.HTMLTemplateFile, ezt)

	} else {
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}

	if ezt.LogLevel >= 3 {
		log.Println("DEBUG-3: Sending", ezt.Type, "Response to", r.RemoteAddr)
	}
}
