// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"github.com/freetaxii/freetaxii-server/lib/systemconfig"
)

// ----------------------------------------------------------------------
// Define Server Type
// ----------------------------------------------------------------------

type ServerType struct {
	SysConfig        *systemconfig.SystemConfigType
	DiscoveryService DiscoveryServiceType
}

type DiscoveryServiceType struct {
	Enabled   bool
	Resources []DiscoveryResourceType
}

type DiscoveryResourceType struct {
	Id          int
	Title       string
	Description string
	Location    string
	Contact     string
	Default     string
	Api_roots   []string
}
