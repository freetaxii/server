// Copyright 2015-2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package headers

import (
	"net/http"
	"strings"

	"github.com/gologme/log"
)

type MediaType struct {
	TAXII21 bool
	TAXII20 bool
	STIX21  bool
	STIX20  bool
	HTML    bool
	JSON    bool
}

func (h *MediaType) ParseTAXII(media string) {
	// If there are spaces after the semicolon, remove all of them
	a := strings.Replace(media, "; ", ";", -1)
	a1 := strings.Split(a, ",")

	for _, v := range a1 {
		if v == "*" || v == "*/*" || v == "application/taxii+json" || v == "application/taxii+json;version=2.1" {
			h.TAXII21 = true
		} else if v == "application/vnd.oasis.taxii+json" || v == "application/vnd.oasis.taxii+json;version=2.0" {
			h.TAXII20 = true
		} else if strings.Contains(v, "application/json") {
			h.JSON = true
		} else if strings.Contains(v, "text/html") {
			h.HTML = true
		}
	}
}

func (h *MediaType) ParseSTIX(media string) {
	a := strings.Replace(media, " ", "", -1)
	a1 := strings.Split(a, ",")

	for _, v := range a1 {
		if v == "*" || v == "*/*" || v == "application/stix+json" || v == "application/stix+json;version=2.1" {
			h.STIX21 = true
		} else if v == "application/vnd.oasis.stix+json" || v == "application/vnd.oasis.stix+json;version=2.0" {
			h.STIX20 = true
		} else if strings.Contains(v, "application/json") {
			h.JSON = true
		} else if strings.Contains(v, "text/html") {
			h.HTML = true
		}
	}
}

// --------------------------------------------------
// Debug HTTP Headers
// --------------------------------------------------

func DebugHttpRequest(r *http.Request) {

	log.Traceln("DEBUG: --------------- BEGIN HTTP DUMP ---------------")
	log.Traceln("DEBUG: Method", r.Method)
	log.Traceln("DEBUG: URL", r.URL)
	log.Traceln("DEBUG: Proto", r.Proto)
	log.Traceln("DEBUG: ProtoMajor", r.ProtoMajor)
	log.Traceln("DEBUG: ProtoMinor", r.ProtoMinor)
	log.Traceln("DEBUG: Header", r.Header)
	log.Traceln("DEBUG: Body", r.Body)
	log.Traceln("DEBUG: ContentLength", r.ContentLength)
	log.Traceln("DEBUG: TransferEncoding", r.TransferEncoding)
	log.Traceln("DEBUG: Close", r.Close)
	log.Traceln("DEBUG: Host", r.Host)
	log.Traceln("DEBUG: Form", r.Form)
	log.Traceln("DEBUG: PostForm", r.PostForm)
	log.Traceln("DEBUG: MultipartForm", r.MultipartForm)
	log.Traceln("DEBUG: Trailer", r.Trailer)
	log.Traceln("DEBUG: RemoteAddr", r.RemoteAddr)
	log.Traceln("DEBUG: RequestURI", r.RequestURI)
	log.Traceln("DEBUG: TLS", r.TLS)
	log.Traceln("DEBUG: --------------- END HTTP DUMP ---------------")
	log.Traceln()
	log.Traceln("DEBUG: --------------- BEGIN HEADER DUMP ---------------")
	for k, v := range r.Header {
		log.Traceln("DEBUG:", k, v)
	}
	log.Traceln("DEBUG: --------------- END HEADER DUMP ---------------")
}
