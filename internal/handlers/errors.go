// Copyright 2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/freetaxii/libstix2/defs"
	"github.com/freetaxii/libstix2/resources/taxiierror"
)

/*
sendUnauthenticatedError - This method will send the correct TAXII error message
for a sessions that is unauthenticated.
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
