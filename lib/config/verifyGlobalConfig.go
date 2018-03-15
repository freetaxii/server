// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package config

import (
	"errors"
	"log"
	"strings"
)

// verifyGlobalDirectives - This method will verify that each required
// configuration directive is present. If an error is found it will be returned.
func (config *ServerConfigType) verifyGlobalConfig() error {
	var problemsFound = 0

	// Protocol Directive
	if config.Global.Protocol != "https" && config.Global.Protocol != "http" {
		log.Infoln("CONFIG: The global.protocol directive must be either https or http")
		problemsFound++
	}

	// TLS Files - only needed if https is defined
	if config.Global.Protocol == "https" {
		problemsFound += config.verifyTLSConfig()
	}

	// Listen Directive
	if config.Global.Listen == "" {
		log.Infoln("CONFIG: The global.listen directive is missing from the configuration file")
		problemsFound++
	}

	// Prefix Directive
	if config.Global.Prefix == "" {
		log.Infoln("CONFIG: The global.prefix directive is missing from the configuration file")
		problemsFound++
	} else {
		if !strings.HasSuffix(config.Global.Prefix, "/") {
			log.Infoln("CONFIG: The global.prefix directive is missing the ending slash '/'")
			problemsFound++
		}
	}

	// Database Configuration Directive
	if config.Global.DbConfig == true && config.Global.DbFile == "" {
		log.Infoln("CONFIG: The global.dbconfig directive is set to true, however, the global.dbfile directive is missing from the configuration file")
		problemsFound++
	}

	// Logging File
	if config.Logging.Enabled == true && config.Logging.LogFile == "" {
		log.Infoln("CONFIG: The logging.logfile directive is missing from the configuration file")
		problemsFound++
	}

	// Return errors if there were any
	if problemsFound > 0 {
		log.Println("ERROR: The Global configuration has", problemsFound, "error(s)")
		return errors.New("ERROR: Configuration errors found")
	}
	return nil
}

/*
verifyTLSConfig - This method will verify that each required TLS configuration
directive are present.
*/
func (config *ServerConfigType) verifyTLSConfig() int {
	var problemsFound = 0

	if config.Global.TLSDir == "" {
		log.Infoln("CONFIG: The global.tlsdir directive is missing from the configuration file")
		problemsFound++
	} else {
		filepath := config.Global.Prefix + config.Global.TLSDir

		if !strings.HasSuffix(config.Global.TLSDir, "/") {
			log.Infoln("CONFIG: The global.tlsdir directive is missing the ending slash '/'")
			problemsFound++
		}

		if !config.exists(filepath) {
			log.Infoln("CONFIG: The TLS path", filepath, "can not be opened")
			problemsFound++
		}
	}

	if config.Global.TLSCrt == "" {
		log.Infoln("CONFIG: The global.tlscrt directive is missing from the configuration file")
		problemsFound++
	} else {
		file := config.Global.Prefix + config.Global.TLSDir + config.Global.TLSCrt
		if !config.exists(file) {
			log.Infoln("CONFIG: The TLS Cert file", file, "can not be opened")
			problemsFound++
		}
	}

	if config.Global.TLSKey == "" {
		log.Infoln("CONFIG: The global.tlskey directive is missing from the configuration file")
		problemsFound++
	} else {
		file := config.Global.Prefix + config.Global.TLSDir + config.Global.TLSKey
		if !config.exists(file) {
			log.Infoln("CONFIG: The TLS Key file", file, "can not be opened")
			problemsFound++
		}
	}

	return problemsFound
}
