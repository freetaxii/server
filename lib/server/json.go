// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"encoding/json"
	"log"
)

// --------------------------------------------------
// Create a TAXII Discovery Response Message
// --------------------------------------------------

// createDiscoveryResponse - This takes in two parameters and will create the TAXII ecoded JSON response
// param: formatPretty - This is a boolean that will tell the Marshal process to format and indent the JSON
// param: index - An integer that lets this method know which discovery service is being handled by this instance
// retval: jsondata - A byte array of JSON encoded data
func (this *ServerType) createDiscoveryResponse(formatPretty bool, index int) []byte {
	var jsondata []byte
	var err error

	tm := this.DiscoveryService.Resources[index]

	if formatPretty == true {
		jsondata, err = json.MarshalIndent(tm, "", "    ")
	} else {
		jsondata, err = json.Marshal(tm)
	}

	if err != nil {
		// If we can not create a status message then there is something
		// wrong with the APIs and nothing is going to work.
		log.Fatal("Unable to create Discovery Response Message")
	}
	return jsondata
}
