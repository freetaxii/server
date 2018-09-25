// Copyright 2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package handlers

import (
	"github.com/freetaxii/libstix2/datastore"
	"github.com/freetaxii/libstix2/objects"
	//"github.com/freetaxii/libstix2/resources"
	"net/http"
	"net/http/httptest"
	"testing"
)

type dummydb struct {
	datastore.Datastorer
}

func (db *dummydb) GetObject(stixid string) (interface{}, error) {
	return objects.NewIndicator("2.0"), nil
}
func (db *dummydb) GetObjectsFromCollection(collectionid string) objects.BundleType {
	b := objects.NewBundle()
	b.SetID("bundle--1234")
	return b
}

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "application/stix+json")

	var objectsSrv ServerHandlerType
	objectsSrv.CollectionID = "1234"
	objectsSrv.DS = &dummydb{}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(objectsSrv.ObjectsServerHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"type":"bundle","id":"bundle--1234","spec_version":"2.0"}` + "\n"

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
