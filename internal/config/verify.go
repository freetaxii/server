// Copyright 2015-2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package config

import (
	"errors"
)

/*
Verify - This method will verify that the configuration file has what it needs.
TODO finish fleshing this out
*/
func (c *ServerConfig) Verify() error {
	var problemsFound = 0

	// --------------------------------------------------
	// Global Configuration
	// --------------------------------------------------
	problemsFound += c.verifyGlobalConfig()

	// --------------------------------------------------
	// Global HTML Configuration
	// --------------------------------------------------
	// If HTML output is turned off globally, then there no need to check the
	// configuration and verify everything is present and valid.
	if c.HTML.Enabled.Value == true {
		problemsFound += c.verifyGlobalHTMLConfig()
	}

	// --------------------------------------------------
	// Discovery Server
	// --------------------------------------------------
	// Only verify the Discovery server configuration if it is enabled.
	if c.DiscoveryServer.Enabled == true {
		problemsFound += c.verifyDiscoveryConfig()

		if c.HTML.Enabled.Value == true {
			problemsFound += c.verifyDiscoveryHTMLConfig()
		}
	}

	// --------------------------------------------------
	// API Root Server
	// --------------------------------------------------
	// Only verify the API Root server configuration if it is enabled.
	if c.APIRootServer.Enabled == true {
		problemsFound += c.verifyAPIRootConfig()

		if c.HTML.Enabled.Value == true {
			problemsFound += c.verifyAPIRootHTMLConfig()
		}
	}

	if problemsFound > 0 {
		c.Logger.Println("ERROR: The configuration has", problemsFound, "error(s)")
		return errors.New("ERROR: Configuration errors found")
	}
	return nil
}
