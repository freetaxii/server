// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"github.com/freetaxii/freetaxii-server/lib/systemconfig"
	"github.com/freetaxii/libtaxii2/objects/api_root"
	"github.com/freetaxii/libtaxii2/objects/discovery"
)

// ----------------------------------------------------------------------
// Define Server Type
// ----------------------------------------------------------------------

type ServerType struct {
	SysConfig        *systemconfig.SystemConfigType
	DiscoveryService DiscoveryServiceType
	ApiRootService   ApiRootServiceType
}

type DiscoveryServiceType struct {
	Enabled   bool
	Resources []discovery.DiscoveryType
}

type ApiRootServiceType struct {
	Enabled   bool
	Resources []api_root.ApiRootType
}
