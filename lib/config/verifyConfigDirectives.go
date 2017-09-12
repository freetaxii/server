// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"errors"
)

// verifyConfigDirectives - This method will verify that each required configuration directive is present
// retval: error
func (this *ServerConfigType) verifyConfigDirectives() error {

	// Listen Directive
	if this.System.Listen == "" {
		return errors.New("The listen directive is missing from the configuration file")
	}

	// Discovery HTML File
	if this.DiscoveryService.HtmlFile == "" {
		return errors.New("The Discovery Service is missing the 'htmlfile' directive in the configuration file")
	}

	// API Root HTML File
	if this.ApiRootService.HtmlFile == "" {
		return errors.New("The API Root Service is missing the 'htmlfile' directive in the configuration file")
	}

	// Discovery Name Directive
	for i, _ := range this.DiscoveryService.Services {
		if this.DiscoveryService.Services[i].Name == "" {
			return errors.New("One or more Discovery Services is missing the 'name' directive in the configuration file")
		} else {
			this.DiscoveryService.Services[i].Path = "/" + this.DiscoveryService.Services[i].Name + "/"
		}
	}

	// API Root Name Directive
	for i, _ := range this.ApiRootService.Services {
		if this.ApiRootService.Services[i].Name == "" {
			return errors.New("One or more API Root Services is missing the 'name' directive in the configuration file")
		} else {
			this.ApiRootService.Services[i].Path = "/" + this.ApiRootService.Services[i].Name + "/"
		}

		if this.ApiRootService.Services[i].Collections.HtmlFile == "" {
			return errors.New("The Collections Service is missing the 'htmlfile' directive in the configuration file")
		}

		this.ApiRootService.Services[i].Collections.Path = this.ApiRootService.Services[i].Path + "collections/"
	}
	return nil
}
