// Copyright 2015-2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file in the root of the source tree.

package config

import (
	"encoding/json"
	"testing"
)

type Dataset struct {
	HTMLConfigType
	Service1 HTMLConfigType
	Service2 HTMLConfigType
	Service3 HTMLConfigType
}

// type HTMLConfigType struct {
// 	Enabled       JSONbool   // User defined in configuration file or set in verifyHTMLConfig()
// 	TemplateDir   JSONstring // User defined in configuration file or set in verifyHTMLConfig()
// 	TemplatePath  JSONstring // Set in verifyHTMLConfig()
// 	TemplateFiles struct {
// 		Discovery   JSONstring // User defined in configuration file or set in verifyHTMLConfig()
// 		APIRoot     JSONstring
// 		Collections JSONstring
// 		Collection  JSONstring
// 		Objects     JSONstring
// 		Manifest    JSONstring
// 	}
// }

var data = `
{
  "Enabled": true,
  "TemplateDir": "html",
  "TemplatePath": "/foo/",
  "TemplateFiles": {
  	"Discovery": "d1",
  	"APIRoot":   "a1"
  	"Collections": "cols1",
  	"Collection": "col1",
  	"Objects": "o1",
  	"Manifest": "m1"
  },
	"Service1": {
		"Enabled": true,
		"TemplateDir": "html",
		"TemplatePath": "/bar/",
		"TemplateFiles": {
	 		"Discovery": "d1",
	  		"APIRoot":   "a2"
	  		"Collections": "cols1",
	  		"Collection": "col1",
	  		"Objects": "o2",
	  		"Manifest": "m1"
	  	}
	},
	"Service2": {
		"Enabled": false,
		"TemplateDir": "html",
		"TemplatePath": "/bar/",
		"TemplateFiles": {
	 		"Discovery": "d1",
	  		"APIRoot":   "a2"
	  		"Collections": "cols1",
	  		"Collection": "col1",
	  		"Objects": "o2",
	  		"Manifest": "m1"
	  	}
	},
	"Service3": {
		"TemplateFiles": {
	 		"Discovery": "d3"
	  	}
	}
}`

// ----------------------------------------------------------------------
func Test_HTMLConfigType(t *testing.T) {
	var c Dataset

	decoder := json.NewDecoder(data)
	err := decoder.Decode(c)

	if err != nil {
		return t.Errorf("error parsing the configuration file: %v", err)
	}

	// t.Log("Test 1: get an error for no collection id")
	// if _, err := sqlGetObjectList(query); err == nil {
	// 	t.Error("no error returned")
	// }

	// t.Log("Test 2: get correct sql statement for object list")
	// query.CollectionID = "aa"
	// testdata = `SELECT t_collection_data.date_added, t_collection_data.stix_id, s_base_object.modified, s_base_object.spec_version FROM t_collection_data JOIN s_base_object ON t_collection_data.stix_id = s_base_object.id WHERE t_collection_data.collection_id = "aa"`
	// if v, _ := sqlGetObjectList(query); testdata != v {
	// 	t.Error("sql statement is not correct")
	// }
}
