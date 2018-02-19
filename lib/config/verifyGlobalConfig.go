// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"errors"
	"log"
	"os"
	"strings"
)

// verifyGlobalDirectives - This method will verify that each required
// configuration directive is present. If an error is found it will be returned.
func (config *ServerConfigType) verifyGlobalConfig() error {
	var problemsFound = 0

	// Protocol Directive
	if config.Global.Protocol != "https" && config.Global.Protocol != "http" {
		log.Println("CONFIG: The global.protocol directive must be either https or http")
		problemsFound++
	}

	// TLS Files - only needed if https is defined
	if config.Global.Protocol == "https" {
		if config.Global.TLSCrt == "" {
			log.Println("CONFIG: The global.tlscrt directive is missing from the configuration file")
			problemsFound++
		} else {
			path := config.Global.Prefix + "etc/tls/"
			problemsFound += config.verifyFileExists(path, config.Global.TLSCrt)
		}

		if config.Global.TLSKey == "" {
			log.Println("CONFIG: The global.tlskey directive is missing from the configuration file")
			problemsFound++
		} else {
			path := config.Global.Prefix + "etc/tls/"
			problemsFound += config.verifyFileExists(path, config.Global.TLSKey)
		}
	}

	// Listen Directive
	if config.Global.Listen == "" {
		log.Println("CONFIG: The global.listen directive is missing from the configuration file")
		problemsFound++
	}

	// Prefix Directive
	if config.Global.Prefix == "" {
		log.Println("CONFIG: The global.prefix directive is missing from the configuration file")
		problemsFound++
	} else {
		if !strings.HasSuffix(config.Global.Prefix, "/") {
			log.Println("CONFIG: The global.prefix directive is missing the ending slash '/'")
			problemsFound++
		}
	}

	// Database Configuration Directive
	if config.Global.DbConfig == true && config.Global.DbFile == "" {
		log.Println("CONFIG: The global.dbconfig directive is set to true, however, the global.dbfile directive is missing from the configuration file")
		problemsFound++
	}

	// HTML Template Dir Directive
	if config.Global.HTMLTemplateDir == "" {
		log.Println("CONFIG: The global.htmltemplatedir directive is missing from the configuration file")
		problemsFound++
	} else {
		if !strings.HasSuffix(config.Global.HTMLTemplateDir, "/") {
			log.Println("CONFIG: The global.htmltemplatedir directive is missing the ending slash '/'")
			problemsFound++
		}
	}

	// Logging File
	if config.Logging.Enabled == true && config.Logging.LogFile == "" {
		log.Println("CONFIG: The logging.logfile directive is missing from the configuration file")
		problemsFound++
	}

	// Return errors if there were any
	if problemsFound > 0 {
		log.Println("ERROR: The Global configuration has", problemsFound, "error(s)")
		return errors.New("ERROR: Configuration errors found")
	}
	return nil
}

// verifyFileExists - This method will check to make sure the file is found on the filesystem
func (config *ServerConfigType) verifyFileExists(path, filename string) int {
	filepath := path + filename
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Println("CONFIG: The TLS file", filename, "can not be opened:", err)
		return 1
	}
	return 0
}
