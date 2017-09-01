// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package systemconfig

import (
	"encoding/json"
	"log"
	"os"
)

// Log Level 1 = basic system logging information, sent to STDOUT unless Enabled = true then it is logged to a file
// Log Level 2 =
// Log Level 3 = detailed debugging information and code troubleshooting (like key variable changes)
// Log Level 4 =
// Log Level 5 = RAW packet/message decode and output

type SystemConfigType struct {
	System struct {
		Listen         string
		Prefix         string
		DbFile         string
		DbFileFullPath string
	}
	Logging struct {
		Enabled         bool
		LogLevel        int
		LogFile         string
		LogFileFullPath string
	}
}

// --------------------------------------------------
// Load Configuration File and Parse JSON
// --------------------------------------------------

func (this *SystemConfigType) LoadSystemConfig(filename string) {

	// Open and read configuration file
	sysConfigFileData, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error opening configuration file: %v", err)
	}

	// --------------------------------------------------
	// Decode JSON configuration file
	// --------------------------------------------------
	// Use decoder instead of unmarshal so we can handle stream data
	decoder := json.NewDecoder(sysConfigFileData)
	err = decoder.Decode(this)

	if err != nil {
		log.Fatalf("error parsing configuration file %v", err)
	}

	// Lets assign the full paths to a few variables so we can use them later
	this.System.DbFileFullPath = this.System.Prefix + "/" + this.System.DbFile
	this.Logging.LogFileFullPath = this.System.Prefix + "/" + this.Logging.LogFile

	if this.Logging.LogLevel >= 5 {
		log.Printf("DEBUG-5: System Configuration Dump %+v\n", this)
	}
}
