// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package headers

import (
	"net/http"

	"github.com/gologme/log"
)

type HttpHeaderType struct {
	DebugLevel int
}

// --------------------------------------------------
// Debug HTTP Headers
// --------------------------------------------------

func (this *HttpHeaderType) DebugHttpRequest(r *http.Request) {

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
