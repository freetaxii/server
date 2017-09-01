// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package headers

import (
	"fmt"
	"net/http"
)

type HttpHeaderType struct {
	DebugLevel int
}

// --------------------------------------------------
// Debug HTTP Headers
// --------------------------------------------------

func (this *HttpHeaderType) DebugHttpRequest(r *http.Request) {

	fmt.Println("DEBUG: --------------- BEGIN HTTP DUMP ---------------")
	fmt.Println("DEBUG: Method", r.Method)
	fmt.Println("DEBUG: URL", r.URL)
	fmt.Println("DEBUG: Proto", r.Proto)
	fmt.Println("DEBUG: ProtoMajor", r.ProtoMajor)
	fmt.Println("DEBUG: ProtoMinor", r.ProtoMinor)
	fmt.Println("DEBUG: Header", r.Header)
	fmt.Println("DEBUG: Body", r.Body)
	fmt.Println("DEBUG: ContentLength", r.ContentLength)
	fmt.Println("DEBUG: TransferEncoding", r.TransferEncoding)
	fmt.Println("DEBUG: Close", r.Close)
	fmt.Println("DEBUG: Host", r.Host)
	fmt.Println("DEBUG: Form", r.Form)
	fmt.Println("DEBUG: PostForm", r.PostForm)
	fmt.Println("DEBUG: MultipartForm", r.MultipartForm)
	fmt.Println("DEBUG: Trailer", r.Trailer)
	fmt.Println("DEBUG: RemoteAddr", r.RemoteAddr)
	fmt.Println("DEBUG: RequestURI", r.RequestURI)
	fmt.Println("DEBUG: TLS", r.TLS)
	fmt.Println("DEBUG: --------------- END HTTP DUMP ---------------")
	fmt.Println("\n")
	fmt.Println("DEBUG: --------------- BEGIN HEADER DUMP ---------------")
	for k, v := range r.Header {
		fmt.Println("DEBUG:", k, v)
	}
	fmt.Println("DEBUG: --------------- END HEADER DUMP ---------------")
}
