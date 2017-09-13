// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

// --------------------------------------------------
// Setup Handler Structs
// --------------------------------------------------
// This struct will handle discovery, api_root, collections, collection, etc
type ServerHandlerType struct {
	Type     string
	Html     bool
	HtmlFile string
	HtmlPath string
	LogLevel int
	Path     string
	Resource interface{}
}
