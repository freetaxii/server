// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package config

import (
	"strings"

	"github.com/gologme/log"
)

/*
verifyGlobalDirectives - This method will verify that each required
configuration directive is present and will return the number of errors found.
*/
func (c *ServerConfigType) verifyGlobalConfig() int {
	var problemsFound = 0

	// Protocol Directive
	if c.Global.Protocol != "https" && c.Global.Protocol != "http" {
		log.Infoln("CONFIG: The global.protocol directive must be either https or http")
		problemsFound++
	}

	// TLS Files - only needed if https is defined
	if c.Global.Protocol == "https" {
		problemsFound += c.verifyTLSConfig()
	}

	// Listen Directive
	if c.Global.Listen == "" {
		log.Infoln("CONFIG: The global.listen directive is missing from the configuration file")
		problemsFound++
	}

	// Prefix Directive
	if c.Global.Prefix == "" {
		log.Infoln("CONFIG: The global.prefix directive is missing from the configuration file")
		problemsFound++
	} else {
		if !strings.HasSuffix(c.Global.Prefix, "/") {
			log.Infoln("CONFIG: The global.prefix directive is missing the ending slash '/'")
			problemsFound++
		}
	}

	// Database Configuration Directive
	if c.Global.DbConfig == true && c.Global.DbFile == "" {
		log.Infoln("CONFIG: The global.dbconfig directive is set to true, however, the global.dbfile directive is missing from the configuration file")
		problemsFound++
	}

	// Logging File
	if c.Logging.Enabled == true && c.Logging.LogFile == "" {
		log.Infoln("CONFIG: The logging.logfile directive is missing from the configuration file")
		problemsFound++
	}

	// ----------------------------------------------------------------------
	// Return number of errors if there are any
	// ----------------------------------------------------------------------
	if problemsFound > 0 {
		log.Println("ERROR: The Global configuration has", problemsFound, "error(s)")
	}
	return problemsFound
}

/*
verifyTLSConfig - This method will verify that each required TLS configuration
directive are present.
*/
func (c *ServerConfigType) verifyTLSConfig() int {
	var problemsFound = 0

	if c.Global.TLSDir == "" {
		log.Infoln("CONFIG: The global.tlsdir directive is missing from the configuration file")
		problemsFound++
	} else {
		filepath := c.Global.Prefix + c.Global.TLSDir

		if !strings.HasSuffix(c.Global.TLSDir, "/") {
			log.Infoln("CONFIG: The global.tlsdir directive is missing the ending slash '/'")
			problemsFound++
		}

		if !c.exists(filepath) {
			log.Infoln("CONFIG: The TLS path", filepath, "can not be opened")
			problemsFound++
		}
	}

	if c.Global.TLSCrt == "" {
		log.Infoln("CONFIG: The global.tlscrt directive is missing from the configuration file")
		problemsFound++
	} else {
		file := c.Global.Prefix + c.Global.TLSDir + c.Global.TLSCrt
		if !c.exists(file) {
			log.Infoln("CONFIG: The TLS Cert file", file, "can not be opened")
			problemsFound++
		}
	}

	if c.Global.TLSKey == "" {
		log.Infoln("CONFIG: The global.tlskey directive is missing from the configuration file")
		problemsFound++
	} else {
		file := c.Global.Prefix + c.Global.TLSDir + c.Global.TLSKey
		if !c.exists(file) {
			log.Infoln("CONFIG: The TLS Key file", file, "can not be opened")
			problemsFound++
		}
	}

	return problemsFound
}
