// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"github.com/freetaxii/freetaxii-server/defs"
	"github.com/freetaxii/freetaxii-server/lib/headers"
	stixObjects "github.com/freetaxii/libstix2/objects"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

// ApiRootServerHandler - This method takes in three parameters. The last parameter
// the index is so that this handler will know which directory service is being called
// in case there is more than one.
// param: w - http.ResponseWriter
// param: r - *http.Request
func (this *ServerHandlerType) ObjectsServerHandler(w http.ResponseWriter, r *http.Request) {
	var mediaType string
	var httpHeaderAccept string
	var jsondata []byte
	var formatpretty bool = false
	var taxiiHeader headers.HttpHeaderType

	// Setup a STIX Bundle to be used for response
	stixBundle := stixObjects.NewBundle()

	// Setup HTML template
	var htmlTemplateResource = template.Must(template.ParseFiles(this.HtmlPath))

	if this.LogLevel >= 3 {
		log.Println("DEBUG-3: Found Request on the", this.Type, "Server Handler from", r.RemoteAddr)
	}

	// We need to put this first so that during debugging we can see problems
	// that will generate errors below.
	if this.LogLevel >= 5 {
		taxiiHeader.DebugHttpRequest(r)
	}

	// If you get a request for a single object, then only send that one object.  Otherwise send them all.

	// This is just sample data
	// TODO move to a database
	i1 := stixObjects.NewIndicator()
	i2 := stixObjects.NewIndicator()

	i1.SetName("Malware C2 Indicator 2016 - File Hash")
	i1.AddLabel("malicious-activity")
	i1.SetPattern("[ file:hashes.'SHA-256' = '4bac27393bdd9777ce02453256c5577cd02275510b2227f473d03f533924f877' ]")
	i1.SetValidFrom(time.Now())
	i1.AddKillChainPhase("lockheed-martin-cyber-kill-chain", "delivery")
	stixBundle.AddObject(i1)

	i2.SetName("Malware C2 Indicator 2016")
	i2.AddLabel("malicious-activity")
	i2.SetPattern("[ ipv4-addr:value = '198.51.100.1/32' ]")
	i2.SetValidFrom(time.Now())
	i2.AddKillChainPhase("lockheed-martin-cyber-kill-chain", "delivery")
	stixBundle.AddObject(i2)

	this.Resource = stixBundle

	w.Header().Add("Strict-Transport-Security", "max-age=86400; includeSubDomains")

	// --------------------------------------------------
	// Decode incoming request message
	// --------------------------------------------------
	httpHeaderAccept = r.Header.Get("Accept")

	if strings.Contains(httpHeaderAccept, defs.TAXII_MEDIA_TYPE) {
		mediaType = defs.TAXII_MEDIA_TYPE + "; " + defs.TAXII_VERSION + "; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		formatpretty = false
		jsondata = this.createTAXIIResponse(formatpretty)
		w.Write(jsondata)
	} else if strings.Contains(httpHeaderAccept, "application/json") {
		mediaType = "application/json; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)
		formatpretty = true
		jsondata = this.createTAXIIResponse(formatpretty)
		w.Write(jsondata)
	} else {
		mediaType = "text/html; charset=utf-8"
		w.Header().Set("Content-Type", mediaType)

		// I needed to convert this to actual JSON since if I just used this.Resource like in other handlers
		// I would get the string output of a Golang struct which is not the same.
		formatpretty = true
		jsondata = this.createTAXIIResponse(formatpretty)
		this.Resource = string(jsondata)
		htmlTemplateResource.ExecuteTemplate(w, this.HtmlFile, this)
	}

	if this.LogLevel >= 3 {
		log.Println("DEBUG-3: Sending", this.Type, "Response to", r.RemoteAddr)
	}
}
