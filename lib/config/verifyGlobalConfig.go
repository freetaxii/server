// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"errors"
	"log"
	"strings"
)

// verifyGlobalDirectives - This method will verify that each required
// configuration directive is present. If an error is found it will be returned.
func (ezt *ServerConfigType) verifyGlobalConfig() error {
	var problemsFound = 0

	// Protocol Directive
	if ezt.Global.Protocol != "https" && ezt.Global.Protocol != "http" {
		log.Println("CONFIG: The global.protocol directive must be either https or http")
		problemsFound++
	}

	// Listen Directive
	if ezt.Global.Listen == "" {
		log.Println("CONFIG: The global.listen directive is missing from the configuration file")
		problemsFound++
	}

	// Prefix Directive
	if ezt.Global.Prefix == "" {
		log.Println("CONFIG: The global.prefix directive is missing from the configuration file")
		problemsFound++
	} else {
		if !strings.HasSuffix(ezt.Global.Prefix, "/") {
			log.Println("CONFIG: The global.prefix directive is missing the ending slash '/'")
			problemsFound++
		}
	}

	// Database Configuration Directive
	if ezt.Global.DbConfig == true && ezt.Global.DbFile == "" {
		log.Println("CONFIG: The global.dbconfig directive is set to true, however, the global.dbfile directive is missing from the configuration file")
		problemsFound++
	}

	// HTML Template Dir Directive
	if ezt.Global.HTMLTemplateDir == "" {
		log.Println("CONFIG: The global.htmltemplatedir directive is missing from the configuration file")
		problemsFound++
	} else {
		if !strings.HasSuffix(ezt.Global.HTMLTemplateDir, "/") {
			log.Println("CONFIG: The global.htmltemplatedir directive is missing the ending slash '/'")
			problemsFound++
		}
	}

	// Logging Config
	if ezt.Logging.Enabled == false && ezt.Logging.LogLevel > 0 {
		log.Println("CONFIG: The logging service is disabled. Setting loglevel to 0.")
		ezt.Logging.LogLevel = 0
	}

	// Return errors if there were any
	if problemsFound > 0 {
		log.Println("ERROR: The Global configuration has", problemsFound, "error(s)")
		return errors.New("ERROR: Configuration errors found")
	}
	return nil
}
