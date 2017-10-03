// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"log"
)

// VerifyServerConfig - This method will verify that the configuration file has what it needs
// TODO finish fleshing this out
func (ezt *ServerConfigType) VerifyServerConfig() error {
	var err error

	// --------------------------------------------------
	// Discovery Server
	// --------------------------------------------------

	err = ezt.verifyGlobalConfig()

	// --------------------------------------------------
	// Discovery Server
	// --------------------------------------------------

	// Only verify the Discovery server configuration if it is enabled.
	if ezt.DiscoveryServer.Enabled == true {
		err = ezt.verifyDiscoveryConfig()
	} else {
		log.Println("CONFIG: The Discovery Server is not enabled in the configuration file")
	}

	if ezt.DiscoveryServer.HTMLEnabled == true {
		err = ezt.verifyDiscoveryHTMLConfig()
	} else {
		log.Println("CONFIG: The Discovery Server is not configured to use HTML output")
	}

	// --------------------------------------------------
	// API Root Server
	// --------------------------------------------------

	// Only verify the API Root server configuration if it is enabled.
	if ezt.APIRootServer.Enabled == true {
		err = ezt.verifyAPIRootConfig()
	} else {
		log.Println("CONFIG: The API Root Server is not enabled in the configuration file")
	}

	if ezt.APIRootServer.HTMLEnabled == true {
		err = ezt.verifyAPIRootHTMLConfig()
	} else {
		log.Println("CONFIG: The API Root Server is not configured to use HTML output")
	}

	if err != nil {
		return err
	}
	return nil
}
