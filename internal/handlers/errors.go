// Copyright 2015-2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/freetaxii/libstix2/defs"
	"github.com/freetaxii/libstix2/resources/taxiierror"
)

/*
sendUnauthenticatedError - This method will send the correct TAXII error message
for a session that is unauthenticated.
*/
func (s *ServerHandler) sendUnauthenticatedError(w http.ResponseWriter) {

	// Setup JSON stream encoder
	j := json.NewEncoder(w)

	w.Header().Set("Content-Type", defs.MEDIA_TYPE_TAXII21)
	w.WriteHeader(http.StatusUnauthorized)

	e := taxiierror.New()
	e.SetTitle("Authentication Required")
	e.SetDescription("The requested resources requires authentication.")
	e.SetErrorCode("401")
	e.SetHTTPStatus("401 Unauthorized")

	j.SetIndent("", "    ")
	j.Encode(e)
}

/*
sendUnsupportedMediaTypeError - This method will send the correct TAXII error
message for a session that requests an unsupported media type in the accept
header.
*/
func (s *ServerHandler) sendUnsupportedMediaTypeError(w http.ResponseWriter) {

	// Setup JSON stream encoder
	j := json.NewEncoder(w)

	w.Header().Set("Content-Type", defs.MEDIA_TYPE_TAXII21)
	w.WriteHeader(http.StatusNotAcceptable)

	e := taxiierror.New()
	e.SetTitle("Wrong Media Type")
	e.SetDescription("The requested media type in the accept header is not supported.")
	e.SetErrorCode("406")
	e.SetHTTPStatus("406 Unsupported Media Type")

	j.SetIndent("", "    ")
	j.Encode(e)
}

/*
sendGetObjectsError - This method will send the correct TAXII error
message for a session that requests some objects but an error is returned.
*/
func (s *ServerHandler) sendGetObjectsError(w http.ResponseWriter) {

	// Setup JSON stream encoder
	j := json.NewEncoder(w)

	w.Header().Set("Content-Type", defs.MEDIA_TYPE_TAXII21)
	w.WriteHeader(http.StatusNotFound)

	e := taxiierror.New()
	e.SetTitle("Get Objects Error")
	e.SetDescription("The request for objects caused an error.")
	e.SetErrorCode("404")
	e.SetHTTPStatus("404 Not Found")

	j.SetIndent("", "    ")
	j.Encode(e)
}
