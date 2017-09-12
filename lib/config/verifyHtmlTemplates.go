// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"errors"
	"os"
)

func (this *ServerConfigType) verifyHtmlTemplateFiles() error {
	// --------------------------------------------------
	// Discovery HTML Resource File Checks
	// --------------------------------------------------
	// Lets check to make sure the Discovery HTML Resource file is on the file system

	filename := this.System.HtmlDir + this.DiscoveryService.HtmlFile
	err := this.verifyHtmlFileExists(filename)
	if err != nil {
		return err
	}

	// Need to check to see if the HTML resource file was redefined at each service level
	for i, _ := range this.DiscoveryService.Services {
		// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
		if this.DiscoveryService.Services[i].HtmlFile == "" {
			this.DiscoveryService.Services[i].HtmlFile = this.DiscoveryService.HtmlFile
		} else {
			// Only test if the file was redefined at this level. No need to retest the inhertied filename since it was already checked
			filename := this.System.HtmlDir + this.DiscoveryService.Services[i].HtmlFile
			err := this.verifyHtmlFileExists(filename)
			if err != nil {
				return err
			}
		}
	}

	// --------------------------------------------------
	// API Root HTML Resource File Checks
	// --------------------------------------------------
	// Lets check to make sure the API Root HTML Resource file is on the file system

	filename2 := this.System.HtmlDir + this.ApiRootService.HtmlFile
	err2 := this.verifyHtmlFileExists(filename2)
	if err2 != nil {
		return err2
	}

	// Need to check to see if the HTML resource file was redefined at each service level
	for i, _ := range this.ApiRootService.Services {
		// If it is not defined at the service level, lets copy in the parent, this will make it easier to work with later on
		if this.ApiRootService.Services[i].HtmlFile == "" {
			this.ApiRootService.Services[i].HtmlFile = this.ApiRootService.HtmlFile
		} else {
			// Only test if the file was redefined at this level. No need to retest the inhertied filename since it was already checked
			filename := this.System.HtmlDir + this.ApiRootService.Services[i].HtmlFile
			err := this.verifyHtmlFileExists(filename)
			if err != nil {
				return err
			}
		}

		// Only test if the file was redefined at this level. No need to reset the inherited filename since it was already checked
		filename := this.System.HtmlDir + this.ApiRootService.Services[i].Collections.HtmlFile
		err := this.verifyHtmlFileExists(filename)
		if err != nil {
			return err
		}

	}

	return nil
}

// verifyHtmlFileExists - This method will check to make sure the HTML resource file is found on the filesystem
// param: file - a string representing the full / relative filename path to the HTML resource file
func (this *ServerConfigType) verifyHtmlFileExists(file string) error {
	filename := this.System.HtmlDir + file
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return errors.New("The following HTML template file is not found on the filesystem:" + filename)
	}
	return nil
}
